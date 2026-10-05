# 09 - Splitting Loan into a package

**New concept:** packages — multiple files, sub-packages, exported vs unexported, `internal/`
**Builds on:** 08

## Task

Move your problem 08 `Loan` into its own package, then use it from `main`.

```
09-packages-modules/
├── main.go                  # package main: Summary(l *loan.Loan) string
├── main_test.go
└── internal/loan/
    ├── errors.go            # given: the sentinels
    ├── validate.go          # TODO: validID, validAmount (unexported)
    ├── validate_test.go     # package loan      — white-box, sees unexported names
    ├── loan.go              # TODO: Loan, New, methods (port from 08)
    └── loan_test.go         # package loan_test — black-box, exported API only
```

- Port 08 into `loan.go`. The constructor is now `loan.New` (not `loan.NewLoan` — see LESSON), and
  there's one new method, `Principal()`.
- `New` and `ApplyPayment` use `validID` / `validAmount` instead of repeating the checks.
- `Summary` in `main.go` formats a loan using only its exported methods:
  `L-100: $600.00 of $1000.00 remaining, payments: 2`. `Summary(nil)` returns `no loan`.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | `validID` / `validAmount` behave like 08's checks | `TestValidID`, `TestValidAmount` |
| 2 | `loan.New` builds a loan, rejects a blank id and a non-positive principal | `TestNew`, `TestNewErrors` |
| 3 | Payments: accepted, rejected with the right sentinel, paid off | `TestApplyPayment` |
| 4 | `Payments()` returns a copy | `TestPayments` |
| 5 | `Summary` formats exactly, including a fresh loan and `nil` | `TestSummary`, `TestSummaryNewLoan`, `TestSummaryNil` |
| 6 | `go doc` shows `New`, `Loan` and its methods, and no fields | `go doc ./problems/09-packages-modules/internal/loan` |
| 7 | `validID` / `validAmount` don't appear in `go doc` | same |

## Run

```bash
go test ./problems/09-packages-modules/...     # runs main, loan, and both test styles
go run ./problems/09-packages-modules
go doc ./problems/09-packages-modules/internal/loan
go doc -all ./problems/09-packages-modules/internal/loan
```

## Stretch (optional)

In `problems/10-interfaces/main.go`, temporarily import `.../09-packages-modules/internal/loan` and run
`go vet ./problems/10-interfaces`. Read the error, then delete the import.
