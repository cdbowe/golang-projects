# Lesson 10 - Interfaces and composition

## Why Go does it this way

- **Interfaces are satisfied implicitly.** A type with the right methods *is* a `Notifier`; there's no `: INotifier`. The type and the interface can live in packages that have never heard of each other.
- **Small interfaces.** The stdlib's most-used interfaces have one method: `io.Reader`, `io.Writer`, `fmt.Stringer`, `error`. "The bigger the interface, the weaker the abstraction."
- **Define interfaces where they're used**, not next to the implementation. The consumer says what it needs; any type that fits can be passed.
- **No inheritance.** Reuse comes from *embedding*: an embedded field's methods are promoted onto the outer type. It's delegation the compiler writes for you, not an is-a relationship.

## Syntax

```go
type Pricer interface {
	Price(sqft float64) float64
}

type FlatRate struct{ PerSqft float64 }

func (f FlatRate) Price(sqft float64) float64 { return f.PerSqft * sqft } // that's it: FlatRate is a Pricer

func Quote(p Pricer, sqft float64) float64 { return p.Price(sqft) } // accept the interface

var _ Pricer = FlatRate{} // compile-time check, costs nothing at runtime

// Type assertion: get the concrete type back
if f, ok := p.(FlatRate); ok { fmt.Println(f.PerSqft) }
```

Embedding:

```go
type Discounted struct {
	Pricer             // embedded interface: Discounted.Price is promoted from it...
	Off    float64
}

func (d Discounted) Price(sqft float64) float64 { // ...unless you define your own, which wins
	return d.Pricer.Price(sqft) * (1 - d.Off)        // reach the embedded one by its type name
}
```

Combining errors (Go 1.20+): `errors.Join(errs...)` returns `nil` if every argument is `nil`; otherwise
an error for which `errors.Is` matches *any* of them.

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| implicit `interface` | `interface` + `: IFoo` | **No direct equivalent.** Nothing declares "implements"; a typo in a method name means "doesn't implement" and you find out at the use site |
| `var _ I = (*T)(nil)` | the `: IFoo` declaration | Opt-in compile-time check, written once, usually next to the type |
| one-method interfaces | `IRepository` with 12 methods | Go prefers many tiny interfaces. A 12-method interface is a smell |
| interface at the consumer | interface next to the implementation | The *caller* package owns the interface. Implementations don't import it |
| embedding | inheritance / base class | No `base.`, no `virtual`/`override`, no polymorphic dispatch to the outer type. Embedded code never calls back into the wrapper |
| `d.Pricer.Price()` | `base.Price()` | You name the embedded field explicitly |
| `p.(FlatRate)` | `p as FlatRate` / pattern match | Use the comma-ok form; the single-result form panics on mismatch |
| nil interface value | `null` | Calling a method on a nil interface panics. A nil *pointer* inside a non-nil interface doesn't — that subtle case is a classic Go bug |
| fakes as plain structs | Moq / NSubstitute | You rarely need a mocking library: a 5-line struct with the right method is the fake |

## Worked example

A small interface, a fake, and a decorator — for appraisals, not notifications:

```go
type Appraiser interface {
	Appraise(address string) (float64, error)
}

// fixedAppraiser is a fake: returns a canned value. No mocking library.
type fixedAppraiser struct{ value float64 }

func (f fixedAppraiser) Appraise(string) (float64, error) { return f.value, nil }

// Capped wraps any Appraiser and limits its result.
type Capped struct {
	Appraiser
	Max float64
}

func (c Capped) Appraise(addr string) (float64, error) {
	v, err := c.Appraiser.Appraise(addr) // delegate
	if err != nil {
		return 0, err
	}
	return min(v, c.Max), nil
}

// Lowest asks every appraiser and keeps the lowest successful value.
func Lowest(as []Appraiser, addr string) (float64, bool) {
	best, found := 0.0, false
	for _, a := range as {
		if v, err := a.Appraise(addr); err == nil && (!found || v < best) {
			best, found = v, true
		}
	}
	return best, found
}
```

`Lowest(as, addr)` accepts `fixedAppraiser`, `Capped`, or anything else with an `Appraise` method.

## Gotchas for C# developers

- Pointer-receiver methods belong to `*T`, not `T`. `var n Notifier = EmailNotifier{}` won't compile; `&EmailNotifier{}` will.
- Embedding a nil interface compiles; calling its promoted method panics. Guard with `if l.Notifier == nil`.
- Don't define an interface "in case" — add it when a second implementation or a test fake needs it.
- Return concrete types, accept interfaces. Returning an interface hides useful methods from callers.
- `any` is `interface{}` — the empty interface every type satisfies. Reaching for it is usually a sign you want generics (12) or a real interface.

## Check yourself

1. `EmailNotifier`'s `Notify` has a pointer receiver. Why does `[]Notifier{email}` work in `main`?
2. With your own `Notify` on `LoggingNotifier` deleted, does the type still satisfy `Notifier`?
3. Why does `TestImplicitSatisfaction` declare its type in the test file?

<details>
<summary>Answers</summary>

1. `email` is `&EmailNotifier{...}`, a `*EmailNotifier`, and that pointer type has the method.
2. Yes: the embedded `Notifier`'s `Notify` is promoted. It compiles and silently skips the logging —
   that's the stretch.
3. To prove the point of implicit interfaces: a type `NotifyAll` has never seen, written somewhere
   else, with no declaration linking it to `Notifier`, still works.

</details>

## Go deeper

- [Tour of Go: interfaces](https://go.dev/tour/methods/9)
- [Effective Go: interfaces](https://go.dev/doc/effective_go#interfaces), [embedding](https://go.dev/doc/effective_go#embedding)
- [Go Proverbs](https://go-proverbs.github.io/) — "The bigger the interface, the weaker the abstraction"
- [pkg.go.dev/errors#Join](https://pkg.go.dev/errors#Join)
