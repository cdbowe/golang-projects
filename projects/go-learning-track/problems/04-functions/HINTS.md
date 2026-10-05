# Hints - 04

<details>
<summary>Hint 1 - direction</summary>

Guards first: bad term or principal returns `(0, 0)`, and a zero rate takes its own branch so you
never divide by zero. Only then the amortization formula from the README.

`TotalFees` is problem 03's loop over a variadic parameter, starting the accumulator at `base`.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- The exponent is negative: `math.Pow(1+r, -float64(termMonths))`. You need `import "math"` — add it
  to the `import` block in `main.go`.
- `termMonths` is an `int`; every arithmetic mix with `float64` needs `float64(termMonths)` explicitly.
- Interest is just `payment * float64(termMonths) - principal` once `payment` is known.
- The results are **named** (`payment`, `totalInterest`), so assign with `=`, never `:=` — see the
  LESSON gotcha about shadowing.
- Inside `TotalFees`, `extra` is a `[]float64`. `for _, f := range extra` is all you need.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
func MonthlyPayment(principal, annualRate float64, termMonths int) (payment, totalInterest float64):
    if termMonths <= 0 || principal <= 0: return 0, 0

    if annualRate == 0:
        payment = principal / float64(termMonths)
        return payment, 0

    r := annualRate / 12
    payment = principal * r / (1 - math.Pow(1+r, -float64(termMonths)))
    totalInterest = payment*float64(termMonths) - principal
    return payment, totalInterest

func TotalFees(base float64, extra ...float64) float64:
    total := base
    for each f in extra: total += f
    return total
```

</details>
