package main

import (
	"cmp"
	"fmt"
	"slices"
	"sync"
)

// Loan is the API's resource.
type Loan struct {
	ID           string  `json:"id"`
	BorrowerName string  `json:"borrower_name"`
	LoanAmount   float64 `json:"loan_amount"`
	Status       string  `json:"status"`
}

// Store is an in-memory loan store. It is given to you complete.
//
// HTTP handlers run on many goroutines at once, so the shared map is guarded
// by a mutex. That's concurrency, which is problem 14 — for now, just call the
// methods. (defer runs a call when the surrounding function returns; here it
// guarantees the unlock.)
type Store struct {
	mu     sync.Mutex
	loans  map[string]Loan
	nextID int
}

// NewStore returns an empty store. IDs start at L-100.
func NewStore() *Store {
	return &Store{loans: make(map[string]Loan), nextID: 100}
}

// Create stores a new loan with status "open" and returns it with its ID.
// It does not validate: that's the handler's job.
func (s *Store) Create(borrower string, amount float64) Loan {
	s.mu.Lock()
	defer s.mu.Unlock()

	l := Loan{
		ID:           fmt.Sprintf("L-%d", s.nextID),
		BorrowerName: borrower,
		LoanAmount:   amount,
		Status:       "open",
	}
	s.nextID++
	s.loans[l.ID] = l
	return l
}

// Get returns the loan with the given ID. ok is false when there is none.
func (s *Store) Get(id string) (l Loan, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, ok = s.loans[id]
	return l, ok
}

// List returns every loan sorted by ID. It never returns nil.
func (s *Store) List() []Loan {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Loan, 0, len(s.loans))
	for _, l := range s.loans {
		out = append(out, l)
	}
	slices.SortFunc(out, func(a, b Loan) int { return cmp.Compare(a.ID, b.ID) })
	return out
}
