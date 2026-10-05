# Hints - 09

<details>
<summary>Hint 1 - direction</summary>

Most of this is a move, not new code. Copy your 08 struct and method bodies into `internal/loan/loan.go`,
put the two checks into `validate.go`, and call those helpers from `New` and `ApplyPayment`. Then write
`Summary` in `main.go` using only methods — fields aren't reachable from `main` anyway.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- `validID` is your 08 `strings.TrimSpace` check returning a `bool`; `validAmount` is `v > 0`.
- `New` → `if !validID(id) { return nil, ErrMissingID }`. Inside package `loan`, the sentinels need no
  prefix; in `main` and `loan_test` they're `loan.ErrMissingID`.
- `Principal()` is new: return the stored principal.
- `Summary`: `fmt.Sprintf` with `%s`, `%.2f`, `%d`. The payment count is `len(l.Payments())`.
- `nil` guard first in `Summary`: `if l == nil { return "no loan" }`.
- The `strings` and `slices` imports move with the code into `loan.go` / `validate.go`.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
// internal/loan/validate.go
validID(id):      TrimSpace(id) != ""
validAmount(v):   v > 0

// internal/loan/loan.go
New(id, principal):
    if !validID(id):            return nil, ErrMissingID
    if !validAmount(principal): return nil, ErrInvalidPrincipal
    return &Loan{...}, nil
ApplyPayment(amount):
    if !validAmount(amount): return ErrInvalidPayment
    ... rest as in 08

// main.go
Summary(l):
    if l == nil: return "no loan"
    return Sprintf("%s: $%.2f of $%.2f remaining, payments: %d",
                   l.ID(), l.Balance(), l.Principal(), len(l.Payments()))
```

</details>
