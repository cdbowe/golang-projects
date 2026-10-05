# Hints - 14

<details>
<summary>Hint 1 - direction</summary>

- `FakeLender.Quote` is one `select`: the delay finishing vs the context ending.
- `AllQuotes` is the LESSON's `AppraiseAll` shape: one result slot per lender, a `WaitGroup`, then a
  sequential pass to split successes from errors.
- `FastestQuote` is a race with bookkeeping: start everyone, then receive results one at a time until
  you get a success, run out of lenders, or the context ends. Make sure losers can always finish.
- `QuoteWithin` is three lines.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- `time.After(f.Delay)` in a `select` with `ctx.Done()`.
- `sync.WaitGroup` with `wg.Go(func() { ... })` (Go 1.25+), then `wg.Wait()`.
- `errors.Join(errs...)` returns `nil` when there's nothing non-nil to join.
- For `FastestQuote`: `context.WithCancel` + `defer cancel()`; a channel **buffered to `len(lenders)`**
  carrying a small `struct{ offer Offer; err error }`; a receive loop that runs at most `len(lenders)`
  times inside a `select` that also watches `ctx.Done()`.
- All-failed error: `errors.Join(append([]error{ErrNoOffers}, errs...)...)`, or
  `fmt.Errorf("%w: %w", ErrNoOffers, errors.Join(errs...))`.
- `context.WithTimeout(ctx, timeout)` + `defer cancel()` for `QuoteWithin`.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
(f *FakeLender) Quote(ctx, amount):
    select:
      case <-time.After(f.Delay): if f.Err != nil: return Offer{}, f.Err; return Offer{f.Name, f.Rate}, nil
      case <-ctx.Done():          return Offer{}, ctx.Err()

AllQuotes(ctx, lenders, amount):
    results := make([]struct{offer; err}, len(lenders))
    wg.Go one func per lender, writing results[i]
    wg.Wait()
    loop results in order: append offer or error
    return offers, errors.Join(errs...)

FastestQuote(ctx, lenders, amount):
    ctx, cancel := WithCancel(ctx); defer cancel()
    ch := make(chan result, len(lenders))
    start one goroutine per lender: ch <- result{lender.Quote(ctx, amount)}
    var errs []error
    for range lenders:
        select:
          case r := <-ch:   if r.err == nil: return r.offer, nil   // deferred cancel stops the rest
                            errs = append(errs, r.err)
          case <-ctx.Done(): return Offer{}, ctx.Err()
    return Offer{}, error matching ErrNoOffers and every errs[i]

QuoteWithin: ctx, cancel := WithTimeout(ctx, timeout); defer cancel(); return FastestQuote(ctx, ...)
```

</details>
