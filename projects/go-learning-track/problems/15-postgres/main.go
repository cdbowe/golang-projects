package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// defaultDSN matches compose.yaml.
const defaultDSN = "postgres://app:secret@localhost:5432/loans?sslmode=disable"

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

// run holds main's logic so that deferred cleanup always runs; log.Fatal in
// main would skip it.
func run(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, cmp.Or(os.Getenv("DATABASE_URL"), defaultDSN))
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer pool.Close()

	if err := Migrate(ctx, pool); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	store := NewStore(pool)
	created, err := store.Create(ctx, "Ana Ruiz", 250000)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	fmt.Printf("created: %+v\n", created)

	loans, err := store.List(ctx)
	if err != nil {
		return fmt.Errorf("list: %w", err)
	}
	fmt.Printf("%d loans in the database\n", len(loans))
	return nil
}
