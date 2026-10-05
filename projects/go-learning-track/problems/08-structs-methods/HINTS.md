# Hints - 08

<details>
<summary>Hint 1 - direction</summary>

Decide what state a loan needs to answer every method: an id, the original principal, and the
payments. Everything else (`Balance`, `IsPaidOff`) can be derived. Then: validate in `NewLoan`,
validate in `ApplyPayment`, and only then mutate.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- Fields: `id string`, `principal float64`, `payments []float64`. Optionally a running `paid float64`.
- `strings.TrimSpace` for the id check — you did this in 05.
- `return &Loan{id: id, principal: principal}, nil` — the nil `payments` slice is ready to append to.
- Every method should use a pointer receiver `(l *Loan)` — `ApplyPayment` must, and the rest match it.
- `slices.Clone` for `Payments()`. Add `"slices"` (and `"strings"`) to the imports.
- `IsPaidOff` with floats: `l.Balance() <= 0` is enough here because overpayment is rejected.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
NewLoan(id, principal):
    if TrimSpace(id) == "": return nil, ErrMissingID
    if principal <= 0:      return nil, ErrInvalidPrincipal
    return &Loan{id, principal}, nil

(l *Loan) Balance():
    principal minus the sum of l.payments

(l *Loan) ApplyPayment(amount):
    if amount <= 0:           return ErrInvalidPayment
    if amount > l.Balance():  return ErrOverpayment
    l.payments = append(l.payments, amount)
    return nil

(l *Loan) IsPaidOff():  l.Balance() <= 0
(l *Loan) Payments():   a clone of l.payments
(l *Loan) ID():         l.id
```

</details>
