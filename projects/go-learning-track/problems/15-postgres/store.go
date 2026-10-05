package main

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("loan not found")
	ErrInvalidStatus = errors.New("invalid status")
)

// Statuses lists every valid Loan.Status.
var Statuses = []string{"open", "closing", "closed"}

// Loan is one row of the loans table.
type Loan struct {
	ID           int64
	BorrowerName string
	LoanAmount   float64
	Status       string
	CreatedAt    time.Time
}

// Store persists loans in Postgres. Safe for concurrent use: the pool hands
// each call its own connection.
type Store struct {
	pool *pgxpool.Pool
}

// NewStore wraps an open connection pool.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// Create inserts a loan and returns it as stored: with its database-assigned
// ID, default status and creation time.
func (s *Store) Create(ctx context.Context, borrower string, amount float64) (Loan, error) {
	// TODO
	return Loan{}, nil
}

// Get returns the loan with the given ID, or ErrNotFound.
func (s *Store) Get(ctx context.Context, id int64) (Loan, error) {
	// TODO
	return Loan{}, nil
}

// List returns every loan ordered by ID.
func (s *Store) List(ctx context.Context) ([]Loan, error) {
	// TODO
	return nil, nil
}

// UpdateStatus sets a loan's status. A status not in Statuses returns
// ErrInvalidStatus without touching the database; an unknown ID returns
// ErrNotFound.
func (s *Store) UpdateStatus(ctx context.Context, id int64, status string) error {
	// TODO
	return nil
}
