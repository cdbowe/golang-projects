# Hints - 07

<details>
<summary>Hint 1 - direction</summary>

- `ApplyPayment`: two guards (nil pointer, non-positive amount), then subtract, then handle "at or
  below zero".
- `WithPayment` already receives a copy — that's what a non-pointer parameter *is*. Change the copy and
  return it.
- `Find`: loop, compare IDs, return the address of the **slice element**, not of your loop variable.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- `l.Balance -= amount` works on a `*Loan` — field access auto-dereferences.
- `if l == nil || amount <= 0 { return }` covers both guards in one line.
- For `Find`, `for i := range loans` gives you the index; `&loans[i]` is the element's address. The
  LESSON's first gotcha shows the version that *compiles but is wrong*.
- `WithPayment` can call `ApplyPayment(&l, amount)` on its own parameter `l`. Its `l` is already a
  copy, so the caller is safe.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
ApplyPayment(l *Loan, amount):
    if l == nil || amount <= 0: return
    l.Balance -= amount
    if l.Balance <= 0:
        l.Balance = 0
        l.Status = "paid"

WithPayment(l Loan, amount) Loan:
    apply the payment to the local copy l   // through its address
    return l

Find(loans []Loan, id) *Loan:
    for i := range loans:
        if loans[i].ID == id: return &loans[i]
    return nil
```

</details>
