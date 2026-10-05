# Lesson 04 - Functions, multiple returns, variadics

## Why Go does it this way

- **Multiple returns are first class.** No `out` parameters, no `Tuple<T1,T2>`, no wrapper struct for two values. This is what makes `(value, error)` the idiom you'll meet in 05.
- **Named results are documentation.** `(payment, totalInterest float64)` tells the caller which is which; godoc shows the names.
- **No overloads, no default arguments.** One name, one signature. Variants get their own name (`Parse`, `ParseInt`, `MustParse`).
- **Variadics are sugar over a slice.** `extra ...float64` *is* a `[]float64` inside the function. Nothing dynamic happens.

## Syntax

```go
func Split(amount float64) (float64, float64) { // two unnamed results
	return amount / 2, amount / 2
}

half, other := Split(100)
half, _ := Split(100)      // _ discards a result you don't want
```

Named results:

```go
func Quote(principal float64, months int) (payment, interest float64) {
	payment = principal / float64(months) // already declared, so = not :=
	interest = 0
	return                                // bare return: returns the named values
}
```

Bare `return` works but is easy to misread in a long function. Returning explicitly is common even
when the results are named: `return payment, interest`.

Shared types collapse:

```go
func MonthlyPayment(principal, annualRate float64, termMonths int) (payment, totalInterest float64)
//                  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ both are float64
```

Variadics:

```go
func Sum(nums ...float64) float64 {
	total := 0.0
	for _, n := range nums { // nums is a []float64
		total += n
	}
	return total
}

Sum()                       // nums is nil, len 0 — no panic
Sum(1, 2, 3)
Sum(values...)              // spread a slice; must be the exact element type
```

Only the **last** parameter can be variadic.

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `(float64, float64)` results | `(double, double)` tuple or `out` params | Not a tuple object — two values on the stack. You can't store them in one variable without a struct |
| named results | XML `<returns>` docs | Go's names are real identifiers, pre-declared to their zero values inside the body |
| bare `return` | — | No C# equivalent. Returns whatever the named results currently hold |
| `extra ...float64` | `params double[] extra` | Nearly identical; spread syntax is postfix `extras...`, not `extras` |
| `_` | discard `_` | Same idea, but in Go it's *required* — an unused real variable won't compile |
| one name, one signature | overloads, optional args | No overloading at all. `Greet()` and `Greet(name)` cannot coexist |
| `float64(months)` | implicit int→double | Go has **no** implicit numeric conversion. `principal / months` with an `int` won't compile |
| `math.Pow(x, y)` | `Math.Pow` | Takes and returns `float64` only. `math.Pow(x, -3)` is fine for negative exponents |

## Worked example

Two results plus a variadic, in a different domain:

```go
// Payoff returns the remaining balance after applying payments, and how many
// of them were actually needed before the balance hit zero.
func Payoff(balance float64, payments ...float64) (remaining float64, used int) {
	remaining = balance // named results start at their zero value

	for _, p := range payments {
		if remaining <= 0 {
			break // stop counting once it's paid off
		}
		remaining -= p
		used++
	}

	if remaining < 0 {
		remaining = 0 // never report a negative balance
	}
	return remaining, used // explicit beats bare when there are guards above
}
```

Calling it:

```go
left, n := Payoff(1000, 400, 400, 400)   // left == 0, n == 3
left, _ = Payoff(1000, 250)              // left == 750, n discarded
sched := []float64{500, 500}
left, n = Payoff(1000, sched...)         // spread
```

## Gotchas for C# developers

- Ignoring a return value is fine for *some* functions, but `go vet` flags a dropped `error`. That lands in 05.
- `payment = ...` inside a function with named results, not `payment := ...` — `:=` would shadow the result and your `return` would hand back the zero value. This is the single most common named-result bug.
- `float64(termMonths)` conversions are mandatory. `1 - math.Pow(...)` is fine because untyped constants adapt; `principal / termMonths` is not.
- Dividing a `float64` by zero gives `+Inf`, not a panic — that's why the zero-rate branch exists. Integer division by zero *does* panic.
- `0.1 + 0.2 != 0.3` here too. Compare money with a tolerance, which is what `closeEnough` in the test file does. (Real money code uses integer cents or a decimal type; `float64` is fine for this exercise.)
- Parameters are passed **by value**, always — including slices and structs (the slice *header* is copied). That's problem 07.

## Check yourself

1. Why does `func Pay(a, b float64)` not mean "a is untyped"?
2. What does `Sum()` do when `nums` is never assigned — nil slice or panic?
3. Inside `func f() (total float64)`, what's wrong with `total := 0.0`?

<details>
<summary>Answers</summary>

1. In a parameter list the type at the end applies to every preceding name in the group: both `a` and
   `b` are `float64`. (C# makes you repeat the type; Go lets you group.)
2. `nums` is a nil `[]float64` with `len(nums) == 0`. `range` over nil runs zero times — no panic.
   Nil slices being usable is a recurring Go theme; more in 06.
3. It declares a **new** local `total` that shadows the result. The compiler accepts it, your bare
   `return` hands back `0`, and the test failure looks like a math bug. Use `total = 0.0`.

</details>

## Go deeper

- [Tour of Go: multiple results](https://go.dev/tour/basics/6)
- [Effective Go: functions](https://go.dev/doc/effective_go#functions)
- [pkg.go.dev/math#Pow](https://pkg.go.dev/math#Pow)
