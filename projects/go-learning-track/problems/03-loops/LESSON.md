# Lesson 03 - for is the only loop

## Why Go does it this way

- **One keyword, four shapes.** No `while`, no `do`, no `foreach`, no `for(;;)` variants to remember. `for` covers all of them.
- **`range` is the safe default.** It hands you index and value, so off-by-one bounds errors mostly disappear.
- **No LINQ.** There is no `.Where().Select()` in the standard library. Go's answer to "transform a collection" is a `for` loop, on purpose — the allocation and the cost are visible in the code.
- **String building is explicit.** `+=` in a loop allocates each time; `strings.Builder` is the `StringBuilder` you already know, and Go makes you reach for it deliberately.

## Syntax

```go
for i := 0; i < 10; i++ { }      // C-style

for cond {  }                    // while-style (no "while" keyword)

for { break }                    // infinite; exit with break or return

for i := range 10 { }            // Go 1.22+: i = 0,1,...,9

for i, v := range things { }     // index + value
for _, v := range things { }     // value only; _ discards
for i := range things { }        // index only
```

`continue` and `break` work as you'd expect. Labels exist for breaking out of nested loops:

```go
outer:
for i := range 3 {
	for j := range 3 {
		if i*j > 2 {
			break outer
		}
	}
}
```

Building a string in a loop:

```go
var b strings.Builder                  // zero value is ready to use, no constructor
for i := range 3 {
	if i > 0 {
		b.WriteString("\n")            // separator before all but the first: no trailing one
	}
	fmt.Fprintf(&b, "line %d", i)      // or b.WriteString(fmt.Sprintf(...))
}
return b.String()
```

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `for i := 0; i < n; i++` | same | No parentheses; braces required |
| `for cond { }` | `while (cond) { }` | `while` is not a keyword in Go |
| `for { }` | `while (true) { }` | Idiomatic, not a smell |
| `for i, v := range xs` | `foreach (var x in xs)` | `range` gives you **two** values; forget the `i` and you'll bind the index into `v` |
| `for i := range 10` | `Enumerable.Range(0, 10)` | Starts at 0 and excludes 10. There is no `..n` inclusive form |
| `strings.Builder` | `StringBuilder` | Zero value works — no `new`. Pass it as `&b` to anything taking an `io.Writer` |
| (nothing) | `do { } while` | Doesn't exist. Use `for { ...; if !cond { break } }` |
| (nothing) | LINQ / `IEnumerable` | No lazy query pipeline in stdlib. Write the loop |
| (nothing) | `foreach` + `yield return` | Generators arrived as range-over-func in Go 1.23; you won't need them this weekend |

## Worked example

Count how many payments clear a threshold, and describe the run — a loop that both aggregates and
builds text:

```go
// LateMonths lists the 1-based months whose payment was short, as "3, 7, 11".
func LateMonths(paid []float64, due float64) string {
	var b strings.Builder
	count := 0

	for i, amount := range paid {
		if amount >= due {
			continue // guard inside the loop keeps the body flat
		}
		if count > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%d", i+1) // range index is 0-based; months are 1-based
		count++
	}

	if count == 0 {
		return "none"
	}
	return b.String()
}
```

Two habits to copy: `continue` as a guard so the interesting branch isn't nested, and `count > 0`
driving the separator so there's no trailing `", "` to trim afterwards.

## Gotchas for C# developers

- `i` declared in `for i := ...` is scoped to the loop. Need it afterwards? Declare it outside.
- `for i := range months` iterates `0..months-1`. The task counts months from **1** — add one, or use
  the C-style form `for m := 1; m <= months; m++`.
- Since Go 1.22 the loop variable is a *new* variable each iteration, so capturing it in a closure is
  safe. Older blog posts warning about this are out of date.
- `s += x` inside a loop is O(n²) in allocations. Fine for 12 lines, wrong for 12,000 — reviewers will
  say `strings.Builder`.
- `range` over a `string` yields **runes with byte offsets**, not chars: the index jumps by 3 on a
  3-byte character. Indexing `s[i]` gives a `byte`. You'll meet this in 06.
- There's no `for ... else`. Track a flag or a counter, as the worked example does.

## Check yourself

1. Write a loop that runs while a balance is above zero, without using `while`.
2. `for i, v := range xs` — what are `i` and `v` for `xs := []string{"a","b"}`?
3. Why does `Schedule` need the separator logic rather than appending `"\n"` after every line?

<details>
<summary>Answers</summary>

1. `for balance > 0 { balance -= payment }` — the condition-only form *is* `while`.
2. `i` is `int` (0 then 1), `v` is `string` ("a" then "b"). `v` is a **copy** of the element; assigning
   to it does not change the slice.
3. Appending after every line leaves a trailing `\n`, and `TestScheduleNoTrailingNewline` checks the
   last byte. Writing the separator *before* every line except the first avoids the trim.

</details>

## Go deeper

- [Tour of Go: for](https://go.dev/tour/flowcontrol/1)
- [Effective Go: for](https://go.dev/doc/effective_go#for)
- [pkg.go.dev/strings#Builder](https://pkg.go.dev/strings#Builder)
