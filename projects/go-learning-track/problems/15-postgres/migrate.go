package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The SQL files are compiled into the binary, so the program runs from any
// directory and never ships without its schema.
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migrate applies every migrations/*.sql file that isn't yet recorded in
// schema_migrations, in file-name order, each in its own transaction. Running
// it again is a no-op. It is given to you complete: you write the SQL files.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version    TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	names, err := fs.Glob(migrationFiles, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	slices.Sort(names)

	for _, name := range names {
		if err := applyMigration(ctx, pool, name); err != nil {
			return err
		}
	}
	return nil
}

func applyMigration(ctx context.Context, pool *pgxpool.Pool, name string) error {
	version := strings.TrimSuffix(path.Base(name), ".sql")

	sql, err := migrationFiles.ReadFile(name)
	if err != nil {
		return fmt.Errorf("read %s: %w", name, err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin %s: %w", version, err)
	}
	defer tx.Rollback(ctx) // no-op once Commit has succeeded

	var applied bool
	err = tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)`, version).Scan(&applied)
	if err != nil {
		return fmt.Errorf("check %s: %w", version, err)
	}
	if applied {
		return nil
	}

	if _, err := tx.Exec(ctx, string(sql)); err != nil {
		return fmt.Errorf("apply %s: %w", version, err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
		return fmt.Errorf("record %s: %w", version, err)
	}
	return tx.Commit(ctx)
}
