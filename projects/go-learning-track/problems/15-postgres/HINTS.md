# Hints - 15

<details>
<summary>Hint 1 - direction</summary>

Schema first: until `001` creates the table, every subtest after the first fails at `reset`. Get
`schema has the expected columns` green, then the constraint subtests, then the `Store` methods.

Each method is one query plus error mapping: `pgx.ErrNoRows` and "zero rows affected" become
`ErrNotFound`; anything else gets wrapped with context.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- DDL: `GENERATED ALWAYS AS IDENTITY PRIMARY KEY`, `CHECK (loan_amount > 0)`, `DEFAULT now()`.
- `002`: `ALTER TABLE loans ADD COLUMN status TEXT NOT NULL DEFAULT 'open' CHECK (status IN (...))`.
- `Create`: `QueryRow(... INSERT ... RETURNING <every column> ...).Scan(...)`. Return columns in
  `Loan`'s field order so `Scan` targets line up.
- `Get`: `QueryRow(...).Scan(...)`; `errors.Is(err, pgx.ErrNoRows)` → `ErrNotFound`.
- `List`: `pool.Query` + `pgx.CollectRows(rows, pgx.RowToStructByPos[Loan])`, selecting columns in
  `Loan`'s field order.
- `UpdateStatus`: `slices.Contains(Statuses, status)`, then `pool.Exec`, then `tag.RowsAffected()`.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```sql
-- 001
CREATE TABLE loans (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    borrower_name TEXT NOT NULL,
    loan_amount   NUMERIC(12,2) NOT NULL CHECK (loan_amount > 0),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- 002: one ALTER TABLE ... ADD COLUMN status ... with NOT NULL, DEFAULT and CHECK
```

```
const loanColumns = "id, borrower_name, loan_amount, status, created_at"   // Loan's field order

Create:  QueryRow("INSERT INTO loans (borrower_name, loan_amount) VALUES ($1, $2) RETURNING " + loanColumns,
                  borrower, amount).Scan(&l.ID, &l.BorrowerName, &l.LoanAmount, &l.Status, &l.CreatedAt)
Get:     QueryRow("SELECT " + loanColumns + " FROM loans WHERE id = $1", id).Scan(...)
         ErrNoRows → ErrNotFound; other err → wrap
List:    rows, _ := Query("SELECT " + loanColumns + " FROM loans ORDER BY id")
         return pgx.CollectRows(rows, pgx.RowToStructByPos[Loan])
UpdateStatus:
         if !slices.Contains(Statuses, status): return ErrInvalidStatus
         tag, err := Exec("UPDATE loans SET status = $1 WHERE id = $2", status, id)
         if tag.RowsAffected() == 0: return ErrNotFound
```

(`loanColumns` is a constant you control, so concatenating it is safe. User input still goes through
`$1`.)

</details>
