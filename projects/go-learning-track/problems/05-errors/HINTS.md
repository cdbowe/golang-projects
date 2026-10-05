# Hints - 05

<details>
<summary>Hint 1 - direction</summary>

Three of the four are a guard and a return: fail → the error value, otherwise `nil`. The fourth calls
the other three in order and adds context to whichever one failed, without losing it.

Whitespace-only names count as missing, so you need to trim before testing for empty.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- [`strings.TrimSpace`](https://pkg.go.dev/strings#TrimSpace) handles `" "`, `"\t\n"`, and `""` in one
  check. Add `"strings"` to the import block.
- Return the sentinel itself: `return ErrMissingBorrower`. Don't build a new error with the same text.
- `RangeError`'s method is on `*RangeError`, so return the address: `return &RangeError{Field: "amount", Value: amount, Min: MinAmount, Max: MaxAmount}`.
- `ValidateTerm` has an `int` input and a `float64` field: `Value: float64(termMonths)`, and `Min: MinTerm` works because `MinTerm` is an untyped constant.
- In `ValidateApplication`, the wrapping verb is `%w` and the prefix string is in the README, exactly:
  `"invalid application: "` — note the colon and the single space.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
func ValidateBorrower(name string) error:
    if strings.TrimSpace(name) == "": return ErrMissingBorrower
    return nil

func ValidateAmount(amount float64) error:
    if amount < MinAmount || amount > MaxAmount:
        return &RangeError{Field: "amount", Value: amount, Min: MinAmount, Max: MaxAmount}
    return nil

func ValidateTerm(termMonths int) error:
    same shape, Field "termMonths", Value float64(termMonths), MinTerm..MaxTerm

func ValidateApplication(borrower string, amount float64, termMonths int) error:
    if err := ValidateBorrower(borrower); err != nil:
        return fmt.Errorf("invalid application: %w", err)
    ... same for ValidateAmount(amount), then ValidateTerm(termMonths)
    return nil
```

Three near-identical `if err := ...; err != nil` blocks is the idiomatic shape here. Collapsing them
needs slices of functions — you'll have the tools for that after 06 and 10.

</details>
