# 05 - Validating loan inputs

**New concept:** `error` as a value — `errors.New`, `fmt.Errorf` + `%w`, `errors.Is` / `errors.As`
**Builds on:** 02, 04

## Task

The sentinel `ErrMissingBorrower`, the `RangeError` type and its `Error()` method are **already
written** for you in `main.go` (structs are problem 08 — here you only return and wrap them).
Implement the four validators.

| Function | Returns on failure | Returns on success |
|---|---|---|
| `ValidateBorrower(name string)` | `ErrMissingBorrower` when `name` is empty or only whitespace | `nil` |
| `ValidateAmount(amount float64)` | `*RangeError{Field: "amount", ...}` outside 1,000-1,000,000 | `nil` |
| `ValidateTerm(termMonths int)` | `*RangeError{Field: "termMonths", ...}` outside 12-480 | `nil` |
| `ValidateApplication(borrower, amount, termMonths)` | the **first** failure, wrapped as `invalid application: <cause>` | `nil` |

Rules:

- Ranges are **inclusive**; the constants are `MinAmount`, `MaxAmount`, `MinTerm`, `MaxTerm`.
- `ValidateApplication` checks borrower → amount → term and stops at the first failure.
- The wrap must preserve the cause: `errors.Is(err, ErrMissingBorrower)` and `errors.As(err, &re)`
  both have to keep working through it. That means `%w`, not `%v`.
- `RangeError.Value` is a `float64`, so the term needs a conversion.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | Empty and whitespace-only names return `ErrMissingBorrower` | `TestValidateBorrowerMissing` |
| 2 | A real name returns `nil` | `TestValidateBorrowerPresent` |
| 3 | Amount bounds are inclusive | `TestValidateAmountInRange` |
| 4 | Out-of-range amount returns `*RangeError` with `Field`, `Value`, `Min`, `Max` filled in | `TestValidateAmountOutOfRange` |
| 5 | Term behaves the same with `Field: "termMonths"` | `TestValidateTerm` |
| 6 | Valid application returns `nil` | `TestValidateApplicationValid` |
| 7 | Wrapped error is still `errors.Is` the sentinel and carries the prefix | `TestValidateApplicationWrapsSentinel` |
| 8 | Wrapped error is still `errors.As` a `*RangeError` | `TestValidateApplicationWrapsRangeError` |
| 9 | Checks run borrower → amount → term, first failure wins | `TestValidateApplicationChecksInOrder` |
| 10 | The range message names the field and the value | `TestRangeErrorMessage` |

## Run

```bash
go test ./problems/05-errors/...
go run ./problems/05-errors
```

## Stretch (optional)

Change one `%w` to `%v` and watch which tests fail. Put it back. That failure signature is worth
recognising on sight.
