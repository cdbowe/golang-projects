# 14 - Fastest rate quote

**New concept:** concurrency — goroutines, channels, `sync.WaitGroup`, `context` cancellation and timeouts
**Builds on:** 05, 10, 11

## Task

Several lenders quote a rate; each is slow. Ask them all at once.

| Function | Behaviour |
|---|---|
| `(*FakeLender).Quote(ctx, amount)` | Wait `Delay`, then return `Offer{Name, Rate}`, or `Err` if set. If `ctx` is done first, return `ctx.Err()` immediately |
| `AllQuotes(ctx, lenders, amount)` | Ask all **concurrently**, wait for all. Return successful offers **in lender order**, and every failure joined (`nil` if none) |
| `FastestQuote(ctx, lenders, amount)` | Ask all concurrently. Return the **first success** and cancel the rest. All failed (or no lenders) → error matching `ErrNoOffers` and each lender's error. `ctx` ends first → error matching `ctx.Err()` |
| `QuoteWithin(ctx, lenders, amount, timeout)` | `FastestQuote` with a deadline of `timeout` from now |

No goroutine may outlive the function that started it. The tests detect leaks.

## How the tests work

Every test runs inside `synctest.Test`, which gives the test a **fake clock**: a 2-second delay
completes instantly, and elapsed times are exact (`took 10ms`, not "about 10ms"). Two failures to
recognise:

| Message | Meaning |
|---|---|
| `took 60ms, want exactly 30ms` | The lenders ran one after another, not concurrently |
| `panic: deadlock: main bubble goroutine has exited but blocked goroutines remain` | A goroutine leaked — usually a send on a channel nobody will ever read |

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | Fake lender: offer or error after its delay; returns at the deadline when the context ends first | `TestFakeLender`, `TestFakeLenderAlreadyCancelled` |
| 2 | `AllQuotes` runs lenders concurrently and keeps lender order | `TestAllQuotesRunsConcurrently` |
| 3 | `AllQuotes` returns the successes plus every error joined | `TestAllQuotesPartialFailure`, `TestAllQuotesEmpty` |
| 4 | `FastestQuote` returns the first success, skipping failures | `TestFastestQuote` |
| 5 | The losers are cancelled the moment there's a winner | `TestFastestQuoteCancelsTheRest` |
| 6 | All fail / no lenders → `ErrNoOffers` plus each cause | `TestFastestQuoteAllFail`, `TestFastestQuoteNoLenders` |
| 7 | `QuoteWithin` answers inside the budget and stops exactly at it | `TestQuoteWithin` |
| 8 | No goroutine leaks | every test (synctest panics on a leak) |

## Run

```bash
go test -v ./problems/14-concurrency/...
go run ./problems/14-concurrency       # real clock: watch the elapsed times
```

`go test -race` is the other standard tool here, but it needs cgo and a C compiler, which this
devcontainer doesn't have (`CGO_ENABLED=0`, no `gcc`).

## Stretch (optional)

Add `CheapestQuote(ctx, lenders, amount)` that waits for everyone and returns the lowest rate. Reuse
`AllQuotes` — it should be a few lines.
