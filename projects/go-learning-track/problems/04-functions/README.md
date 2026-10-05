# 04 - Monthly payment calculator

**New concept:** multiple return values, named results, variadic parameters
**Builds on:** 02, 03

## Task

`MonthlyPayment(principal, annualRate float64, termMonths int) (payment, totalInterest float64)`

Standard amortization. `annualRate` is a decimal fraction (`0.065` == 6.5%).

```
r        = annualRate / 12
payment  = principal * r / (1 - (1+r)^-termMonths)
interest = payment * termMonths - principal
```

| Condition | Result |
|---|---|
| `termMonths <= 0` | `(0, 0)` — checked first, at any rate |
| `principal <= 0` | `(0, 0)` — checked first, at any rate |
| `annualRate == 0` | `payment = principal / termMonths`, `interest = 0` |

Do not round. The tests compare to the cent.

`TotalFees(base float64, extra ...float64) float64` returns `base` plus every extra. No extras
returns `base`.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | 250,000 @ 6.5% / 360mo → payment 1580.17, interest 318,861.22 | `TestMonthlyPayment30Year` |
| 2 | 30,000 @ 4.99% / 60mo → 566.00 / 3,959.97 | `TestMonthlyPayment5YearAuto` |
| 3 | 500,000 @ 7.25% / 180mo → 4,564.31 / 321,576.59 | `TestMonthlyPayment15Year` |
| 4 | Zero rate splits principal evenly, no division by zero | `TestMonthlyPaymentZeroRate` |
| 5 | Bad term or principal returns `(0, 0)`, at any rate including zero | `TestMonthlyPaymentGuards`, `TestMonthlyPaymentGuardsZeroRate` |
| 6 | `TotalFees` works with zero, two, and five extras | `TestTotalFeesNoExtras`, `TestTotalFeesWithExtras` |
| 7 | A `[]float64` spreads into `TotalFees` with `...` | `TestTotalFeesSpreadSlice` |

## Run

```bash
go test ./problems/04-functions/...
go run ./problems/04-functions
```

## Stretch (optional)

Switch `MonthlyPayment` between explicit `return payment, interest` and bare `return` with named
results. Keep whichever reads better to you — then say why in one line.
