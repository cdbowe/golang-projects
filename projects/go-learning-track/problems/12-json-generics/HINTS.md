# Hints - 12

<details>
<summary>Hint 1 - direction</summary>

- Tags first: once they're right, most of `ParseLoans` and `EncodeLoan` is one library call each.
- `ParseLoans` takes an `io.Reader`, and you need to reject unknown keys — both point to a `Decoder`
  rather than `json.Unmarshal`.
- The generic functions are problem 03/06 loops with type parameters. `Number` is the only part that's
  genuinely new: it lists types, it doesn't declare methods.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- Tag format: `` `json:"key"` ``, `` `json:"key,omitempty"` ``, `` `json:"-"` ``.
- [`json.NewDecoder`](https://pkg.go.dev/encoding/json#NewDecoder), `DisallowUnknownFields`, `Decode(&loans)`.
- [`json.Marshal`](https://pkg.go.dev/encoding/json#Marshal) for `EncodeLoan`.
- Wrap: `fmt.Errorf("parse loans: %w", err)`.
- `Number`: a union of types with `~` in front of each, e.g. `~int | ~float64`. Add `~int64` if you like.
- `Filter`/`Map`: build a **new** slice with `append` (or `make([]U, 0, len(items))` for `Map`).

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
type Loan struct {
    ID string `json:"id"`
    ... one tag per field; Documents gets omitempty; InternalNotes gets "-"
}

ParseLoans(r):
    dec := json.NewDecoder(r)
    dec.DisallowUnknownFields()
    var loans []Loan
    if err := dec.Decode(&loans); err != nil: return nil, Errorf("parse loans: %w", err)
    return loans, nil

EncodeLoan(l):  return json.Marshal(l)

type Number interface { ~int | ~int64 | ~float64 }

Filter[T](items, keep):
    var out []T
    for each it: if keep(it): out = append(out, it)
    return out

Map[T, U](items, f):   out := make([]U, 0, len(items)); append f(it) for each; return out
Sum[N](nums):          total += n for each; return total
```

</details>
