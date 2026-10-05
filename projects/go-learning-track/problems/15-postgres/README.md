# 15 - Loans in Postgres

**New concept:** databases — `pgx`, plain-SQL migrations, Testcontainers-go integration tests
**Builds on:** 05, 11, 13, 14

## Task

Write the two migrations and the four `Store` methods. `Migrate` (in `migrate.go`, given complete)
applies `migrations/*.sql` in name order and records each in `schema_migrations`.

**`001_create_loans.sql`** — create `loans`, columns in this order:

| Column | Type | Rules |
|---|---|---|
| `id` | `BIGINT` | `GENERATED ALWAYS AS IDENTITY PRIMARY KEY` |
| `borrower_name` | `TEXT` | `NOT NULL` |
| `loan_amount` | `NUMERIC(12,2)` | `NOT NULL`, must be `> 0` |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL`, defaults to `now()` |

**`002_add_loan_status.sql`** — add `status TEXT`, `NOT NULL`, default `'open'`, allowed values
`open`, `closing`, `closed` only.

**`store.go`:**

| Method | Behaviour |
|---|---|
| `Create(ctx, borrower, amount)` | Insert, return the row **as stored** (ID, default status, `created_at`) — one round trip |
| `Get(ctx, id)` | The loan, or `ErrNotFound` |
| `List(ctx)` | Every loan, ordered by `id` |
| `UpdateStatus(ctx, id, status)` | Status not in `Statuses` → `ErrInvalidStatus` (no query). Unknown ID → `ErrNotFound` |

All queries use `$1`-style parameters. Never build SQL with string concatenation.

## Docker

The test starts a throwaway `postgres:17-alpine` container through Testcontainers. The devcontainer
runs its own Docker daemon (the `docker-in-docker` feature), so this works after a container rebuild.
Without Docker:

| Option | Command |
|---|---|
| Skip the integration test | `go test -short ./...` |
| Use a Postgres you already have | `TEST_DATABASE_URL=postgres://... go test ./problems/15-postgres/...` (tables are truncated) |

For `go run`, start `compose.yaml`'s database first. Its default address (`localhost:5432`) is what
`main.go` uses when `DATABASE_URL` is unset.

## Acceptance criteria

All checks are subtests of `TestPostgres`.

| # | Criterion | Checked by |
|---|---|---|
| 1 | `loans` has exactly the five columns, in order, with the types above | `schema has the expected columns` |
| 2 | Both migrations recorded; re-running `Migrate` is a no-op | `migrations are recorded...` |
| 3 | `Create` returns ID, default status, DB-set `created_at` | `create returns the stored row` |
| 4 | `NUMERIC` keeps cents exactly | `amounts keep their cents` |
| 5 | `Get` round-trips every field; unknown ID → `ErrNotFound` | `get returns what create stored`, `get unknown id is ErrNotFound` |
| 6 | `List` ordered by ID; empty table → no rows, no error | `list is ordered by id`, `list of an empty table` |
| 7 | `UpdateStatus` updates, and rejects bad statuses and unknown IDs | `update status`, `update status rejects...` |
| 8 | The schema itself rejects zero/negative amounts, missing borrower, bad or null status | `constraints reject bad rows` |

## Run

```bash
go test -v ./problems/15-postgres/...                       # needs Docker (or TEST_DATABASE_URL)
go test -v -run 'TestPostgres/update' ./problems/15-postgres/...
docker compose -f problems/15-postgres/compose.yaml up -d --wait   # for go run
go run ./problems/15-postgres
docker compose -f problems/15-postgres/compose.yaml down          # when finished
```

## Stretch (optional)

Add `003_add_loans_status_index.sql` creating an index on `status`, and a subtest checking that
`schema_migrations` now lists three versions.
