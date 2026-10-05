# Lesson 06 - Slices and maps

## Why Go does it this way

- **A slice is a 3-word view** (pointer, `len`, `cap`) onto a backing array. Cheap to pass, but two slices can share memory.
- **`append` may or may not reallocate.** Always assign the result: `s = append(s, x)`.
- **Maps are built in**, not a library type. Iteration order is **randomized on purpose**, so nobody depends on it.
- **Nil is usable for reads.** `len(nilSlice) == 0`, `range nilMap` runs zero times, `nilMap[k]` returns the zero value. Only *writing* to a nil map panics.

## Syntax

```go
var s []int                 // nil slice: len 0, cap 0, appendable
s = append(s, 1, 2, 3)      // MUST reassign
s2 := make([]int, 0, 10)    // len 0, cap 10: preallocate when you know the size
s3 := []int{4, 5, 6}        // literal
part := s3[1:3]             // [5 6] — shares memory with s3
dup := slices.Clone(s3)     // independent copy

m := map[string]int{}       // empty, writable
m2 := make(map[string]int)  // same thing
var m3 map[string]int       // nil: readable, NOT writable
m["a"] = 1
delete(m, "a")

v, ok := m["a"]             // comma-ok: ok == false when key absent
if _, ok := m["b"]; !ok { } // existence check only

for i, v := range s3 { }    // index, value (copy)
for k, v := range m { }     // random order
```

Useful `slices` functions (Go 1.21+): `Sort`, `Clone`, `Equal`, `Contains`, `Index`, `Max`, `Sorted`.

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `[]T` slice | `List<T>` / `Span<T>` | A *view*, not an owner. Sub-slicing shares memory, like `Span<T>`, but it's heap-safe and everywhere |
| `append(s, x)` | `list.Add(x)` | Returns a new header. Forget `s =` and the item is lost |
| `make([]T, 0, n)` | `new List<T>(n)` | Same idea: capacity hint |
| `map[K]V` | `Dictionary<K,V>` | Missing key returns the zero value — **no `KeyNotFoundException`** |
| `v, ok := m[k]` | `TryGetValue(k, out v)` | Same shape; built into the language, not a method |
| map iteration | `Dictionary` iteration | Go randomizes order every run. C# looks ordered until it isn't |
| `slices.Clone(s)` | `list.ToList()` | Explicit copy — required whenever you sort or mutate data you don't own |
| `slices.Sort(s)` | `list.Sort()` | **In place**, like C#. `OrderBy` has no direct equivalent — clone first |
| nil slice | `null` list | A nil slice is *safe*: `len`, `range`, `append` all work on it |

## Worked example

Count loan statuses, then report them in a stable order:

```go
// StatusCounts tallies how many loans are in each status.
func StatusCounts(statuses []string) map[string]int {
	counts := make(map[string]int, len(statuses)) // capacity hint, not a limit
	for _, s := range statuses {
		counts[s]++ // missing key reads as 0, so no existence check needed
	}
	return counts
}

// Busiest returns the status with the most loans; ok is false for no input.
func Busiest(counts map[string]int) (status string, ok bool) {
	best := -1
	for s, n := range counts {
		if n > best || (n == best && s < status) { // tie-break: random order must not leak out
			status, best = s, n
		}
	}
	return status, best >= 0
}
```

Two habits worth copying: `counts[s]++` uses the zero value instead of `if !ok { counts[s] = 0 }`, and
the tie-break keeps the result deterministic despite random iteration.

## Gotchas for C# developers

- `s[1:3]` then `append` to it can **overwrite** `s[3]` in the original. Clone before you mutate data you don't own.
- `var m map[string]int; m["x"] = 1` panics: `assignment to entry in nil map`. Use `make` or `{}`.
- `for _, v := range s { v *= 2 }` changes nothing — `v` is a copy. Use `s[i] *= 2`.
- `slices.Sort(prices)` sorts the **caller's** slice. Function arguments don't protect you here.
- Slices aren't comparable with `==` (except to `nil`). Use `slices.Equal`.
- `len` is a built-in function, not a property: `len(s)`, never `s.Length` or `s.Count`.
- Nil vs empty slice: both have `len 0`; prefer checking `len(s) == 0` over `s == nil`.

## Check yourself

1. `a := []int{1,2,3}; b := a[:2]; b = append(b, 99)` — what is `a` now?
2. `m := map[string][]float64{}; m["x"] = append(m["x"], 5)` — does it work on the first call?
3. Why does `TestSortedZips` call `SortedZips` 20 times?

<details>
<summary>Answers</summary>

1. `[1 2 99]`. `b` had `cap 3`, so `append` wrote into the shared array instead of reallocating.
2. Yes. `m["x"]` is a nil slice (the zero value), and appending to nil allocates. That's the grouping
   idiom — no "create the list if missing" step like C#'s `TryGetValue` + `Add`.
3. Map order is randomized on every `range`. An unsorted result can come out sorted once by luck; it
   won't do it 20 times in a row.

</details>

## Go deeper

- [Go blog: Slices intro](https://go.dev/blog/slices-intro) — the backing-array diagrams are the best explanation
- [Tour of Go: slices](https://go.dev/tour/moretypes/7), [maps](https://go.dev/tour/moretypes/19)
- [pkg.go.dev/slices](https://pkg.go.dev/slices), [pkg.go.dev/maps](https://pkg.go.dev/maps)
