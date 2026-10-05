# 12 - Loans from JSON, generic helpers

**New concept:** `encoding/json` with struct tags, and generics (`[T any]`, constraints)
**Builds on:** 06, 08, 10, 11

## Task

**Part 1 — JSON.** Tag `Loan`'s fields so they map to these keys, then implement `ParseLoans` and
`EncodeLoan`.

| Field | JSON key | Rule |
|---|---|---|
| `ID` | `id` | |
| `BorrowerName` | `borrower_name` | |
| `LoanAmount` | `loan_amount` | |
| `Status` | `status` | |
| `Documents` | `documents` | left out of the output when empty |
| `InternalNotes` | — | never read or written |

- `ParseLoans(r io.Reader)` decodes a JSON **array** of loans. Unknown keys are an error. Every error
  starts with `parse loans: ` and wraps the cause with `%w`.
- `EncodeLoan(l)` returns compact JSON, fields in struct order.

**Part 2 — generics.** Implement `Filter[T]`, `Map[T, U]` and `Sum[N Number]`, and define the `Number`
constraint. `Filter` must not modify its input. `Sum` must accept named types such as `type cents int`.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | `testdata/loans.json` decodes into 3 fully populated loans | `TestParseLoansFile` |
| 2 | `[]` decodes to zero loans, no error | `TestParseLoansEmptyArray` |
| 3 | Truncated, unknown-key, wrong-type, non-array and empty input all fail, prefixed and wrapped | `TestParseLoansErrors` |
| 4 | Exact compact JSON; empty `documents` omitted; `InternalNotes` never written | `TestEncodeLoan` |
| 5 | Encode → parse gives back the same loan, minus notes | `TestRoundTrip` |
| 6 | `Filter` keeps order, handles no matches, leaves input untouched | `TestFilter` |
| 7 | `Map` works across types, including a stdlib func as `f` | `TestMap` |
| 8 | `Sum` works for `int`, `float64`, and `cents` | `TestSum` |

## Run

```bash
go test -v ./problems/12-json-generics/...
go run ./problems/12-json-generics
```

## Stretch (optional)

Write `MaxBy[T any, K cmp.Ordered](items []T, key func(T) K) (T, bool)` and use it in `main` to print
the largest loan.
