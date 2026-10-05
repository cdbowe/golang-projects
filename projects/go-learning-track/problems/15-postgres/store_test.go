package main

import (
	"context"
	"errors"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// newTestPool returns a migrated connection pool. By default it starts a
// throwaway Postgres container (needs Docker). Set TEST_DATABASE_URL to use an
// existing database instead — its tables are truncated between subtests.
func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test: needs Postgres; run without -short")
	}
	ctx := context.Background()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		ctr, err := postgres.Run(ctx, "postgres:17-alpine",
			postgres.WithDatabase("loans"),
			postgres.WithUsername("app"),
			postgres.WithPassword("secret"),
			postgres.BasicWaitStrategies(),
		)
		testcontainers.CleanupContainer(t, ctr)
		if err != nil {
			t.Fatalf("start postgres container (is Docker running?): %v", err)
		}
		if dsn, err = ctr.ConnectionString(ctx, "sslmode=disable"); err != nil {
			t.Fatalf("connection string: %v", err)
		}
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	return pool
}

// reset empties the loans table and restarts its ID sequence.
func reset(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), `TRUNCATE loans RESTART IDENTITY`); err != nil {
		t.Fatalf("reset: %v (did migration 001 create the loans table?)", err)
	}
}

func mustCreate(t *testing.T, s *Store, borrower string, amount float64) Loan {
	t.Helper()
	l, err := s.Create(context.Background(), borrower, amount)
	if err != nil {
		t.Fatalf("Create(%q, %v): %v", borrower, amount, err)
	}
	return l
}

// One container for the whole suite; every subtest starts from an empty table.
func TestPostgres(t *testing.T) {
	pool := newTestPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	t.Run("schema has the expected columns", func(t *testing.T) {
		rows, _ := pool.Query(ctx, `
			SELECT a.attname || ' ' || format_type(a.atttypid, a.atttypmod)
			FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
			WHERE c.relname = 'loans' AND a.attnum > 0 AND NOT a.attisdropped
			ORDER BY a.attnum`)
		got, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			t.Fatal(err)
		}
		want := []string{
			"id bigint",
			"borrower_name text",
			"loan_amount numeric(12,2)",
			"created_at timestamp with time zone",
			"status text",
		}
		if !slices.Equal(got, want) {
			t.Errorf("loans columns =\n %q\nwant\n %q\n(001 creates four, 002 adds status)", got, want)
		}
	})

	t.Run("migrations are recorded and re-running is a no-op", func(t *testing.T) {
		if err := Migrate(ctx, pool); err != nil {
			t.Fatalf("second Migrate: %v", err)
		}
		rows, _ := pool.Query(ctx, `SELECT version FROM schema_migrations ORDER BY version`)
		versions, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{"001_create_loans", "002_add_loan_status"}; !slices.Equal(versions, want) {
			t.Errorf("schema_migrations = %v, want %v", versions, want)
		}
	})

	t.Run("create returns the stored row", func(t *testing.T) {
		reset(t, pool)
		before := time.Now().Add(-time.Minute)

		l := mustCreate(t, store, "Ana Ruiz", 250000)

		if l.ID != 1 || l.BorrowerName != "Ana Ruiz" || l.LoanAmount != 250000 || l.Status != "open" {
			t.Errorf("created = %+v, want ID 1, Ana Ruiz, 250000, open", l)
		}
		if l.CreatedAt.Before(before) {
			t.Errorf("CreatedAt = %v, want it set by the database to about now", l.CreatedAt)
		}
	})

	t.Run("amounts keep their cents", func(t *testing.T) {
		reset(t, pool)
		l := mustCreate(t, store, "Ben", 1234.56)

		got, err := store.Get(ctx, l.ID)
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if got.LoanAmount != 1234.56 {
			t.Errorf("LoanAmount = %v, want 1234.56", got.LoanAmount)
		}
	})

	t.Run("get returns what create stored", func(t *testing.T) {
		reset(t, pool)
		created := mustCreate(t, store, "Cy", 95000)

		got, err := store.Get(ctx, created.ID)
		if err != nil {
			t.Fatalf("Get(%d): %v", created.ID, err)
		}
		if got.ID != created.ID || got.BorrowerName != created.BorrowerName ||
			got.LoanAmount != created.LoanAmount || got.Status != created.Status ||
			!got.CreatedAt.Equal(created.CreatedAt) {
			t.Errorf("Get = %+v, want %+v", got, created)
		}
	})

	t.Run("get unknown id is ErrNotFound", func(t *testing.T) {
		reset(t, pool)
		if _, err := store.Get(ctx, 999); !errors.Is(err, ErrNotFound) {
			t.Errorf("Get(999) err = %v, want ErrNotFound", err)
		}
	})

	t.Run("list is ordered by id", func(t *testing.T) {
		reset(t, pool)
		for _, name := range []string{"A", "B", "C"} {
			mustCreate(t, store, name, 1000)
		}

		loans, err := store.List(ctx)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		var names []string
		for _, l := range loans {
			names = append(names, l.BorrowerName)
		}
		if want := []string{"A", "B", "C"}; !slices.Equal(names, want) {
			t.Errorf("List borrowers = %v, want %v", names, want)
		}
	})

	t.Run("list of an empty table", func(t *testing.T) {
		reset(t, pool)
		loans, err := store.List(ctx)
		if err != nil || len(loans) != 0 {
			t.Errorf("List = (%v, %v), want (empty, nil)", loans, err)
		}
	})

	t.Run("update status", func(t *testing.T) {
		reset(t, pool)
		l := mustCreate(t, store, "Dee", 5000)

		if err := store.UpdateStatus(ctx, l.ID, "closing"); err != nil {
			t.Fatalf("UpdateStatus: %v", err)
		}
		got, err := store.Get(ctx, l.ID)
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if got.Status != "closing" {
			t.Errorf("Status = %q, want %q", got.Status, "closing")
		}
	})

	t.Run("update status rejects unknown statuses and ids", func(t *testing.T) {
		reset(t, pool)
		l := mustCreate(t, store, "Eve", 5000)

		if err := store.UpdateStatus(ctx, l.ID, "funded"); !errors.Is(err, ErrInvalidStatus) {
			t.Errorf("UpdateStatus(funded) err = %v, want ErrInvalidStatus", err)
		}
		if err := store.UpdateStatus(ctx, 999, "closed"); !errors.Is(err, ErrNotFound) {
			t.Errorf("UpdateStatus(id 999) err = %v, want ErrNotFound", err)
		}
		if got, _ := store.Get(ctx, l.ID); got.Status != "open" {
			t.Errorf("Status = %q after rejected updates, want %q", got.Status, "open")
		}
	})

	// The database protects itself even from code that skips the Go checks.
	t.Run("constraints reject bad rows", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{"zero amount", `INSERT INTO loans (borrower_name, loan_amount) VALUES ('X', 0)`},
			{"negative amount", `INSERT INTO loans (borrower_name, loan_amount) VALUES ('X', -1)`},
			{"missing borrower", `INSERT INTO loans (loan_amount) VALUES (1000)`},
			{"unknown status", `INSERT INTO loans (borrower_name, loan_amount, status) VALUES ('X', 1000, 'funded')`},
			{"null status", `INSERT INTO loans (borrower_name, loan_amount, status) VALUES ('X', 1000, NULL)`},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				reset(t, pool)
				if _, err := pool.Exec(ctx, tc.sql); err == nil {
					t.Errorf("insert succeeded, want a constraint violation:\n%s", tc.sql)
				}
			})
		}
	})
}
