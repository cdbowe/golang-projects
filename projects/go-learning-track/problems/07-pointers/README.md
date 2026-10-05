# 07 - Mutate vs copy

**New concept:** pointers — `&`, `*`, pass-by-value, when to use a pointer
**Builds on:** 05, 06

## Task

`Loan` is given (structs are 08; here you only read and write fields). Implement three functions.

| Function | Behaviour |
|---|---|
| `ApplyPayment(l *Loan, amount float64)` | Changes the caller's loan in place |
| `WithPayment(l Loan, amount float64) Loan` | Returns an updated copy; the caller's loan is untouched |
| `Find(loans []Loan, id string) *Loan` | Pointer to the matching element **inside the slice**, or `nil` |

Payment rules, the same for `ApplyPayment` and `WithPayment`:

| Condition | Result |
|---|---|
| `l == nil` (pointer version) | no-op, no panic |
| `amount <= 0` | no change |
| payment brings balance to 0 or below | `Balance = 0`, `Status = "paid"` |
| otherwise | `Balance -= amount` |

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | `ApplyPayment` changes the caller's loan | `TestApplyPaymentMutates` |
| 2 | Exact payoff and overpayment both end at `0` / `paid` | `TestApplyPaymentPaysOff`, `TestApplyPaymentOverpayFloorsAtZero` |
| 3 | Zero and negative payments are ignored | `TestApplyPaymentIgnoresNonPositive` |
| 4 | `nil` loan doesn't panic | `TestApplyPaymentNil` |
| 5 | `WithPayment` returns the update and leaves the original alone, all fields carried over | `TestWithPaymentLeavesOriginal`, `TestWithPaymentPaysOff` |
| 6 | `Find` returns `&loans[i]`, so writes reach the slice | `TestFindReturnsSliceElement` |
| 7 | `Find` returns `nil` for a missing ID or nil slice | `TestFindMissing` |

## Run

```bash
go test ./problems/07-pointers/...
go run ./problems/07-pointers
```

## Stretch (optional)

If `WithPayment` repeats the payment rules, rewrite it in two lines by reusing `ApplyPayment`. Explain
in one sentence why that doesn't change the caller's loan.
