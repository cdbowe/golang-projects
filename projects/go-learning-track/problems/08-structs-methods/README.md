# 08 - The Loan type

**New concept:** structs + methods — value vs pointer receivers, `NewX` constructors
**Builds on:** 05, 06, 07

## Task

Design the `Loan` struct's fields yourself, then implement its constructor and five methods. The
error sentinels are given.

| Member | Behaviour |
|---|---|
| `NewLoan(id, principal) (*Loan, error)` | Blank/whitespace `id` → `ErrMissingID`. `principal <= 0` → `ErrInvalidPrincipal`. On error the loan is `nil` |
| `ID() string` | The id |
| `Balance() float64` | Principal minus the accepted payments |
| `ApplyPayment(amount) error` | `amount <= 0` → `ErrInvalidPayment`. `amount > Balance()` → `ErrOverpayment`. Rejected payments change nothing. Paying the exact balance is allowed |
| `IsPaidOff() bool` | `true` once the balance is zero |
| `Payments() []float64` | Accepted payments in order, as a **copy** |

Fields stay unexported: callers only go through `NewLoan` and the methods.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | A new loan has its id, full balance, no payments, isn't paid off | `TestNewLoan` |
| 2 | Bad id / principal → matching sentinel and a `nil` loan | `TestNewLoanMissingID`, `TestNewLoanInvalidPrincipal` |
| 3 | A payment persists on the loan the caller holds | `TestApplyPaymentPersists` |
| 4 | Non-positive and over-balance payments are rejected and leave no trace | `TestApplyPaymentRejectsNonPositive`, `TestApplyPaymentRejectsOverpayment` |
| 5 | Paying the exact balance pays it off | `TestPaidOff` |
| 6 | History holds accepted payments only, in order | `TestPaymentsHistory` |
| 7 | Editing `Payments()`'s result doesn't edit the loan | `TestPaymentsReturnsCopy` |
| 8 | Two loans share no state | `TestLoansAreIndependent` |
| 9 | Fields are unexported | review |

## Run

```bash
go test ./problems/08-structs-methods/...
go run ./problems/08-structs-methods
```

## Stretch (optional)

Add `PercentPaid() float64` (0-100). Pick its receiver deliberately, and say in one line why it should
match the others even though it doesn't mutate.
