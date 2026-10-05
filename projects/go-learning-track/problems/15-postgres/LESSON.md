# Lesson 15 - Postgres with pgx

## Why Go does it this way

- **SQL stays SQL.** The Go norm is hand-written queries plus a thin driver, not an ORM. `pgx` is the Postgres driver most Go services use; `database/sql` is the generic stdlib interface over many drivers.
- **A pool, not connections.** `pgxpool.Pool` is safe for concurrent use; each call borrows a connection and returns it. Create one at startup, pass it down, close it at shutdown.
- **`context` on every call.** Queries are cancelled when the request (or the test) is — the 14 machinery, applied to I/O.
- **Migrations are files.** Numbered `.sql` files applied in order and recorded in a table. Tools like `goose` and `golang-migrate` do exactly what `migrate.go` does, with more features.

## Syntax

```go
pool, err := pgxpool.New(ctx, "postgres://user:pw@host:5432/db?sslmode=disable")
defer pool.Close()

// No result rows
tag, err := pool.Exec(ctx, `UPDATE appraisals SET value = $1 WHERE id = $2`, 510000, id)
tag.RowsAffected() // 0 → nothing matched

// Exactly one row
var v float64
err = pool.QueryRow(ctx, `SELECT value FROM appraisals WHERE id = $1`, id).Scan(&v)
if errors.Is(err, pgx.ErrNoRows) { /* not found */ }

// Insert and read back in one round trip
err = pool.QueryRow(ctx,
	`INSERT INTO appraisals (address, value) VALUES ($1, $2) RETURNING id, ordered_at`,
	addr, value).Scan(&a.ID, &a.OrderedAt)

// Many rows → a slice
rows, _ := pool.Query(ctx, `SELECT id, address, value FROM appraisals ORDER BY id`)
list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Appraisal]) // fields in SELECT order

// Transaction
tx, err := pool.Begin(ctx)
defer tx.Rollback(ctx)  // no-op after a successful Commit
// ... tx.Exec / tx.QueryRow ...
err = tx.Commit(ctx)
```

`//go:embed migrations/*.sql` (in `migrate.go`) compiles files into the binary as an `embed.FS`.

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `pgx` | Npgsql | Same role. `pgx` has its own API; `database/sql` is the ADO.NET-style abstraction |
| `pgx.CollectRows` + `RowToStructByPos` | Dapper `Query<T>` | Maps by column position (or by name with `` `db:"..."` `` tags and `RowToStructByName`) |
| `pgxpool.Pool` | `NpgsqlDataSource` | Long-lived and shared. Don't open a pool per request |
| `$1, $2` | `@p1` parameters | Postgres-native placeholders, positional |
| `.sql` migration files | EF Core migrations | No model diffing, no generated C#. You write the DDL |
| `pgx.ErrNoRows` | `SingleOrDefault()` returning `null` | It's an **error** from `Scan`; map it to your domain's `ErrNotFound` |
| `defer tx.Rollback(ctx)` | `using var tx` | Same safety net: any early return rolls back |
| Testcontainers-go | Testcontainers for .NET | Same project, same model; cleanup via `t.Cleanup` |
| `testing.Short()` | `[Trait("Category","Integration")]` filter | `go test -short` skips tests that call `t.Skip` when `Short()` is true |

## Worked example

A store for a different table, showing the error mapping:

```go
var ErrNoAppraisal = errors.New("appraisal not found")

type Appraisal struct {
	ID        int64
	Address   string
	Value     float64
	OrderedAt time.Time
}

func (s *AppraisalStore) Get(ctx context.Context, id int64) (Appraisal, error) {
	var a Appraisal
	err := s.pool.QueryRow(ctx,
		`SELECT id, address, value, ordered_at FROM appraisals WHERE id = $1`, id,
	).Scan(&a.ID, &a.Address, &a.Value, &a.OrderedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Appraisal{}, ErrNoAppraisal // callers shouldn't need to import pgx
	}
	if err != nil {
		return Appraisal{}, fmt.Errorf("get appraisal %d: %w", id, err)
	}
	return a, nil
}

func (s *AppraisalStore) Cancel(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM appraisals WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("cancel appraisal %d: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNoAppraisal
	}
	return nil
}
```

## Gotchas for C# developers

- `fmt.Sprintf("... WHERE id = %d", id)` into SQL is injection-shaped. Always pass parameters.
- `Scan` targets are positional: the order must match the `SELECT` list, and the count must match exactly.
- A `Query` you don't read to the end holds its connection. `CollectRows` closes for you; manual loops must `defer rows.Close()` and check `rows.Err()`.
- `NUMERIC` into `float64` works but can round very large or very precise values. Real money code uses integer cents or a decimal type.
- `TIMESTAMPTZ` comes back in the session's time zone; compare with `time.Time.Equal`, not `==`.
- Constraints belong in the schema too. Go validation protects the happy path; `CHECK` and `NOT NULL` protect the data from every other writer.

## Check yourself

1. Why does `Create` use `INSERT ... RETURNING` rather than an insert followed by a `SELECT`?
2. `UpdateStatus` on a missing ID: what does Postgres return, and how do you notice?
3. Why does `migrate.go` run each file inside a transaction?

<details>
<summary>Answers</summary>

1. One round trip, and no race: the row you get back is the one you inserted, with database-set
   values (`id`, `status` default, `created_at`) filled in.
2. No error — the `UPDATE` succeeds and touches zero rows. Check `tag.RowsAffected() == 0`.
3. A file that fails halfway is rolled back entirely, and isn't recorded as applied. Postgres DDL is
   transactional, so the schema is never left half-migrated.

</details>

## Go deeper

- [pkg.go.dev/github.com/jackc/pgx/v5](https://pkg.go.dev/github.com/jackc/pgx/v5) — `CollectRows`, `RowTo*`
- [Go docs: Accessing relational databases](https://go.dev/doc/database/)
- [Testcontainers for Go: Postgres module](https://golang.testcontainers.org/modules/postgres/)
- [pkg.go.dev/embed](https://pkg.go.dev/embed)
