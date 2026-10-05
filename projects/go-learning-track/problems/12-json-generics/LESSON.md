# Lesson 12 - JSON and generics

## Why Go does it this way

- **Struct tags are string metadata** read by reflection at runtime: `` `json:"loan_amount"` ``. One field can carry tags for several libraries (`json`, `db`, `validate`).
- **Only exported fields are serialized.** Lowercase fields are invisible to `encoding/json`, tag or no tag.
- **Generics are deliberately small.** Type parameters on functions and types, constraints as interfaces, inference at call sites. No generic *methods*, no specialization, no variance.
- **Constraints say which operators work.** `[T any]` allows only what every type supports (assign, pass around). To use `+` or `<`, the constraint must list types that support it.

## Syntax

JSON:

```go
type Offer struct {
	Lender string   `json:"lender"`
	Rate   float64  `json:"rate"`
	Fees   []string `json:"fees,omitempty"` // dropped when nil or empty
	Secret string   `json:"-"`              // never serialized
	note   string                           // unexported: ignored entirely
}

data, err := json.Marshal(offer)            // []byte, compact
err = json.Unmarshal(data, &offer)          // decode into a pointer

dec := json.NewDecoder(r)                   // streaming, from any io.Reader
dec.DisallowUnknownFields()                 // unknown key → error
err = dec.Decode(&offers)
```

Generics:

```go
func Count[T any](items []T, match func(T) bool) int { ... } // T: any type

type Ordered interface {                     // a constraint is an interface...
	~int | ~float64 | ~string                // ...listing a type set. ~int = int or any type built on int
}

func Largest[T Ordered](a, b T) T {
	if a > b { return a }                    // > is allowed: every type in the set supports it
	return b
}

Largest(3, 7)                                // T inferred as int
Largest[float64](3, 7.5)                     // or explicit
```

A function value is an ordinary value: pass a named function (`strconv.Itoa`) or a literal
(`func(n int) bool { return n > 2 }`).

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `` `json:"loan_amount"` `` | `[JsonPropertyName("loan_amount")]` | A raw string, not a typed attribute. A typo like `` `json:loan_amount` `` (no quotes) is silently ignored; `go vet` catches some |
| `json:"-"` | `[JsonIgnore]` | Same |
| `omitempty` / `omitzero` (1.24+) | `JsonIgnoreCondition.WhenWritingDefault` | `omitempty` drops empty slices/maps/strings and zeros; `omitzero` drops only the zero value |
| `DisallowUnknownFields()` | `JsonUnmappedMemberHandling.Disallow` | Opt-in, on a `Decoder` only — `json.Unmarshal` can't do it |
| case-insensitive key match | `PropertyNameCaseInsensitive = true` | **On by default** in Go: `"ID"`, `"id"`, `"Id"` all fill `ID` |
| `[T any]` | `<T>` | Type parameters in square brackets |
| constraint interface | `where T : INumber<T>` | A type *set* (`~int \| ~float64`), not a method requirement |
| `func(T) bool` | `Func<T, bool>` / `Predicate<T>` | Plain function type; no delegate declaration |
| `Filter` / `Map` | LINQ `Where` / `Select` | Not in the stdlib; eager, return slices, no deferred execution |
| no generic methods | generic methods on classes | `func (s *Store) Get[T]()` doesn't compile. Use a generic function |

## Worked example

```go
type Quote struct {
	Lender string  `json:"lender"`
	Rate   float64 `json:"rate"`
}

// Index builds a map keyed by any comparable key.
func Index[T any, K comparable](items []T, key func(T) K) map[K]T {
	m := make(map[K]T, len(items))
	for _, it := range items {
		m[key(it)] = it
	}
	return m
}

func LoadQuotes(data []byte) (map[string]Quote, error) {
	var qs []Quote
	if err := json.Unmarshal(data, &qs); err != nil {
		return nil, fmt.Errorf("load quotes: %w", err)
	}
	return Index(qs, func(q Quote) string { return q.Lender }), nil
}
```

`comparable` is a built-in constraint: types usable as map keys and with `==`.

## Gotchas for C# developers

- Decode into a **pointer**: `json.Unmarshal(data, loans)` without `&` returns an error (`non-pointer`).
- JSON numbers decode into `float64` when the target is `any`. Use typed structs.
- A nil slice encodes as `null`; an empty slice as `[]`. API consumers notice.
- `int | float64` without `~` rejects `type cents int`. The error says so: `possibly missing ~ for int`.
- `Filter` that reuses the input (`out := items[:0]`) is a known trick — and it overwrites the caller's slice. Allocate a new slice.
- Don't reach for generics first. Write it concretely; generalize when the second type shows up.

## Check yourself

1. Without tags, which of `Loan`'s fields would `testdata/loans.json` fill?
2. Why does `Sum` fail to compile with `type Number interface{}` but `Filter` compiles with `any`?
3. What's the difference between `omitempty` on `[]string{}` and on `nil`?

<details>
<summary>Answers</summary>

1. `ID` and `Status`. Matching is case-insensitive, so `id` → `ID` and `status` → `Status`. Nothing maps
   `borrower_name` to `BorrowerName` — the underscore doesn't match.
2. `Sum` uses `+`, and an empty constraint allows any type, including ones without `+`. `Filter` only
   stores and returns values, which every type supports.
3. None — both are omitted, because `omitempty` checks `len == 0`. Without `omitempty`, nil writes
   `null` and empty writes `[]`.

</details>

## Go deeper

- [Go blog: JSON and Go](https://go.dev/blog/json)
- [pkg.go.dev/encoding/json](https://pkg.go.dev/encoding/json) — the tag options list
- [Tutorial: Getting started with generics](https://go.dev/doc/tutorial/generics)
- [Go blog: When To Use Generics](https://go.dev/blog/when-generics)
