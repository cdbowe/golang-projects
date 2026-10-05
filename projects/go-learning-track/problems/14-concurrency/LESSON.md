# Lesson 14 - Goroutines, channels, context

## Why Go does it this way

- **Goroutines are cheap.** A few KB of stack, scheduled by the runtime onto OS threads. Starting 10,000 is normal. There's no `async`/`await` and no function colouring: any function can run concurrently with `go f()`.
- **"Share memory by communicating."** Prefer passing values over a channel to locking shared data. Mutexes exist (see 13's `Store`) for when shared state is simpler.
- **Cancellation is explicit and cooperative.** A `context.Context` is passed as the first parameter; callees watch `ctx.Done()` and give up. Nothing is ever killed from outside.
- **Every goroutine needs an exit plan.** A goroutine blocked forever is a leak — the runtime won't collect it.

## Syntax

```go
go work()                       // start; no handle, no result

ch := make(chan Offer)          // unbuffered: send blocks until someone receives
ch := make(chan Offer, 3)       // buffered: send blocks only when 3 are waiting
ch <- o                         // send
o := <-ch                       // receive
close(ch)                       // no more sends; receivers drain then get zero values

var wg sync.WaitGroup
wg.Go(func() { work() })        // Go 1.25+: Add(1) + go + Done() in one call
wg.Wait()                       // block until every wg.Go func has returned

select {                        // wait on whichever is ready first
case o := <-results:
	use(o)
case <-ctx.Done():
	return ctx.Err()            // context.Canceled or context.DeadlineExceeded
case <-time.After(time.Second): // a one-off timer
}

ctx, cancel := context.WithCancel(parent)       // cancel() ends ctx and its children
ctx, cancel := context.WithTimeout(parent, d)   // ends by itself after d
defer cancel()                                   // ALWAYS: releases the timer and children
```

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `go f()` | `Task.Run(f)` | Returns nothing. Results and errors come back through channels or shared state |
| goroutine | `Task` | No `await`; blocking a goroutine is fine and normal |
| `chan T` | `Channel<T>` | Built into the language; unbuffered by default (a rendezvous) |
| `select` | `Task.WhenAny` | Waits on channel operations, not tasks. `default:` makes it non-blocking |
| `sync.WaitGroup` | `Task.WhenAll` | Counts goroutines; carries no results |
| `context.Context` | `CancellationToken` | Also carries deadlines (and request-scoped values). Passed explicitly, first parameter |
| `ctx.Done()` | `token.Register` / `IsCancellationRequested` | A channel, so it composes with `select` |
| `context.WithTimeout` | `CancellationTokenSource.CancelAfter` | `cancel` must be called even if the timeout fires |
| `errors.Join` | `AggregateException` | Just an error; `errors.Is` matches any member |
| `sync.Mutex` | `lock` | Not re-entrant. Copying a struct that holds one is a bug (`go vet` flags it) |

## Worked example

Check several appraisers concurrently; collect every result in input order:

```go
type result struct {
	value float64
	err   error
}

func AppraiseAll(ctx context.Context, as []Appraiser, addr string) []result {
	results := make([]result, len(as)) // one slot per goroutine: no shared writes, no lock
	var wg sync.WaitGroup
	for i, a := range as {
		wg.Go(func() {
			v, err := a.Appraise(ctx, addr) // i and a are per-iteration (Go 1.22+)
			results[i] = result{v, err}
		})
	}
	wg.Wait()
	return results
}
```

And a first-answer race, with a buffered channel so losers never block:

```go
func FirstResponder(ctx context.Context, svcs []Service) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // stops the losers when we return

	answers := make(chan string, len(svcs)) // room for everyone: nobody blocks on send
	for _, s := range svcs {
		go func() {
			if a, err := s.Ask(ctx); err == nil {
				answers <- a
			}
		}()
	}
	select {
	case a := <-answers:
		return a, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
```

(This one hangs if every service fails before `ctx` ends — your `FastestQuote` must count failures.)

## Gotchas for C# developers

- A send on an unbuffered channel blocks until someone receives. Return early and the sender waits forever: a leak.
- `defer cancel()` right after every `WithCancel`/`WithTimeout`. `go vet` warns when you don't.
- Writing to the same slice **index** from many goroutines is a race; writing to **different** indexes is fine.
- `main` returning ends the program, running goroutines included. No implicit join.
- A panic in any goroutine crashes the whole program. There's no unobserved-task exception.
- `time.Sleep` in production code to "wait for goroutines" is a bug. Use `WaitGroup` or channels.

## Check yourself

1. In `FirstResponder`, what goes wrong if `answers` is unbuffered?
2. Why can `AppraiseAll` write `results[i]` from many goroutines without a mutex?
3. `ctx.Err()` after `WithTimeout` fires vs after you call `cancel()` — what's the difference?

<details>
<summary>Answers</summary>

1. Only the first answer is received. Every other service that succeeds blocks on `answers <- a`
   forever, after the function has returned — one leaked goroutine per extra success.
2. Each goroutine owns a different element; nothing is read or written by two of them. `wg.Wait()`
   then guarantees every write happened before `results` is returned.
3. `context.DeadlineExceeded` when the timer fired first; `context.Canceled` when `cancel()` ran
   first. Callers use `errors.Is` to tell "too slow" from "abandoned".

</details>

## Go deeper

- [Go blog: Go Concurrency Patterns: Context](https://go.dev/blog/context)
- [Go blog: Pipelines and cancellation](https://go.dev/blog/pipelines)
- [pkg.go.dev/sync#WaitGroup.Go](https://pkg.go.dev/sync#WaitGroup.Go), [pkg.go.dev/testing/synctest](https://pkg.go.dev/testing/synctest)
- [Effective Go: concurrency](https://go.dev/doc/effective_go#concurrency)
