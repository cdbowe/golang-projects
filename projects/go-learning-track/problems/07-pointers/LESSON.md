# Lesson 07 - Pointers and pass-by-value

## Why Go does it this way

- **Everything is passed by value.** Ints, strings, structs, slices, maps — the callee always gets a copy. What differs is *what* gets copied: a struct copies all its fields; a slice copies a 3-word header that still points at the same array.
- **Pointers make sharing explicit.** If a function takes `*Loan`, it can change your loan. If it takes `Loan`, it can't. The signature tells you.
- **No pointer arithmetic.** Pointers are safe references, not C pointers. The garbage collector handles lifetime; returning `&local` is fine.
- **No class/struct split.** C# decides reference-vs-value by *type* (`class` vs `struct`). Go decides at each *use*: `Loan` or `*Loan`.

## Syntax

```go
// Asterisk (*) operator == declare pointer var, de-reference pointer
// Ampersand (&) operator == returns address/reference to store in a pointer var

x := 10
p := &x      // p is *int: the address of x
*p = 20      // write through the pointer: x is now 20
fmt.Println(*p, x) // 20 20

var q *int   // nil pointer
if q != nil { fmt.Println(*q) } // *q on nil panics: always guard

type Point struct{ X, Y int }
pt := Point{1, 2}
pp := &pt
pp.X = 9     // auto-deref for fields: same as (*pp).X = 9

func bump(n *int) { *n++ }
bump(&x)     // caller passes the address
```

Pointer into a slice element:

```go
nums := []int{1, 2, 3}
second := &nums[1] // address of the element itself
*second = 99       // nums is now [1 99 3]
```

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `func f(l Loan)` | passing a `struct` | Same copy semantics — but in Go this applies to *every* struct, including your "entity" types |
| `func f(l *Loan)` | passing a `class` instance | You choose per parameter, not per type. Nothing is a reference type "by nature" |
| `&x` | `ref x` | `ref` is call-scoped; a Go pointer is a value you can store, return, and keep |
| `*p = v` | assigning through `ref` | Explicit dereference; fields auto-deref (`p.X`), plain values don't |
| `nil` pointer | `null` reference | Same `NullReferenceException` risk — Go calls it `nil pointer dereference` and panics |
| return `&local` | — | Legal and safe in Go: the compiler moves `local` to the heap (escape analysis) |
| `*Loan` as "maybe a loan" | `Loan?` / `null` | Returning `nil` for "not found" is idiomatic. Comma-ok is the other common choice |
| no `readonly` param | `in` parameter | Go has no const-pointer. If you hand out a pointer, the receiver can write through it |

## Worked example

Pointers to plain values — a rate adjustment applied in place vs as a copy:

```go
// Bump raises *rate by delta, capped at max. A nil rate is ignored.
func Bump(rate *float64, delta, max float64) {
	if rate == nil {
		return
	}
	*rate += delta
	if *rate > max {
		*rate = max
	}
}

// Bumped returns the raised rate and leaves the caller's value alone.
func Bumped(rate, delta, max float64) float64 {
	return min(rate+delta, max)
}

// Highest returns a pointer to the largest rate in rates, or nil if empty.
func Highest(rates []float64) *float64 {
	var best *float64
	for i := range rates {
		if best == nil || rates[i] > *best {
			best = &rates[i] // index, not the range value: &v would be a copy
		}
	}
	return best
}
```

```go
r := 0.065
Bump(&r, 0.01, 0.07)          // r == 0.07
n := Bumped(r, 0.01, 0.08)    // n == 0.08, r still 0.07
```

## Gotchas for C# developers

- `for _, l := range loans { if l.ID == id { return &l } }` compiles, and returns the address of a **copy**. Writes go nowhere. Use the index.
- A pointer into a slice can go stale: a later `append` that reallocates leaves your pointer aimed at the old array. Don't hold `&s[i]` across appends.
- `p.Field` auto-derefs, so `p` looks like a C# reference. `p == nil` still needs guarding before you touch a field.
- Passing a pointer for "performance" on small structs is usually slower (heap escape + GC). Default to values; use pointers when you need to share mutation, or the struct is large.
- `*` in a type (`*Loan`) means "pointer to"; `*` in an expression (`*p`) means "dereference". Same symbol, two jobs.
- There's no `new Loan { ... }` requirement — `&Loan{ID: "L1"}` is the idiom. `new(Loan)` exists but is rarely used for structs.

## Check yourself

1. `func reset(s []int) { s[0] = 0 }` — does the caller see the change? What about `func grow(s []int) { s = append(s, 1) }`?
2. Why does `TestFindReturnsSliceElement` compare `p != &loans[1]`, not just the balance?
3. When would you return `Loan` instead of `*Loan` from a function?

<details>
<summary>Answers</summary>

1. `reset`: yes — the header copy points at the same array. `grow`: no — `append` changed the callee's
   header only. That's why `append` returns a slice you must assign.
2. Pointer identity is the actual contract: "this *is* the element". A balance check could pass by
   accident if `Find` wrote through some other route; the identity check can't.
3. When the caller should get an independent snapshot (like `WithPayment`), or the value is small and
   immutable in practice. Return a pointer when the caller should share and mutate one instance.

</details>

## Go deeper

- [Tour of Go: pointers](https://go.dev/tour/moretypes/1)
- [Go FAQ: when are function parameters passed by value?](https://go.dev/doc/faq#pass_by_value)
- [Go FAQ: should I define methods on values or pointers?](https://go.dev/doc/faq#methods_on_values_or_pointers) — read again after 08
