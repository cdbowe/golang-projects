package main

import "fmt"

// Loan is given to you — structs and methods are problem 08. Here you only
// read and write its fields, through a value or through a pointer.
type Loan struct {
	ID      string
	Balance float64
	Status  string // "active" or "paid"
}

// ApplyPayment reduces the loan's balance by amount, in place. A nil loan or a
// non-positive amount changes nothing. The balance never goes below 0; when it
// reaches 0 the status becomes "paid".
func ApplyPayment(l *Loan, amount float64) {
	if l == nil || amount < 0.0 {
		return
	}

	l.Balance = max(0.0, l.Balance-amount)
	if l.Balance == 0.0 {
		l.Status = "paid"
	}
}

// WithPayment returns a copy of l with the payment applied, using the same
// rules as ApplyPayment. The caller's loan is not changed.
func WithPayment(l Loan, amount float64) Loan {
	ApplyPayment(&l, amount) // params are passed by value by default, so function receives a copy of `l` straightaway.
	return l
}

// Find returns a pointer to the loan in loans with the given ID, or nil when
// there is none. The pointer refers to the element inside the slice, so a
// write through it changes loans.
func Find(loans []Loan, id string) *Loan {
	for i := range loans {
		if loans[i].ID == id {
			return &loans[i]
		}
	}

	return nil
}

func main() {
	loans := []Loan{
		{ID: "L-100", Balance: 1200, Status: "active"},
		{ID: "L-200", Balance: 300, Status: "active"},
	}

	preview := WithPayment(loans[0], 200)
	fmt.Printf("preview %s: %.2f (stored still %.2f)\n", preview.ID, preview.Balance, loans[0].Balance)

	if l := Find(loans, "L-200"); l != nil {
		ApplyPayment(l, 300)
	}
	fmt.Printf("stored %s: %.2f %s\n", loans[1].ID, loans[1].Balance, loans[1].Status)
}
