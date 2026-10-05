package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrMissingID        = errors.New("loan id is required")
	ErrInvalidPrincipal = errors.New("principal must be positive")
	ErrInvalidPayment   = errors.New("payment must be positive")
	ErrOverpayment      = errors.New("payment exceeds balance")
)

// Loan tracks a principal and the payments made against it.
type Loan struct {
	id        string
	principal float64
	balance   float64
	payments  []float64
}

// NewLoan validates its inputs and returns a ready-to-use loan.
// See README.md for the error each bad input returns.
func NewLoan(id string, principal float64) (*Loan, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrMissingID
	}

	if principal <= 0 {
		return nil, ErrInvalidPrincipal
	}

	return &Loan{id: id, principal: principal, balance: principal}, nil
}

// ID returns the loan's identifier.
func (l *Loan) ID() string {
	return l.id // fields of a struct auto-deref
}

// Balance returns the current saved balance.
func (l *Loan) Balance() float64 {
	return l.balance
}

// Principal returns the original principal amount.
func (l *Loan) Principal() float64 {
	return l.principal
}

// Balance returns the principal minus everything paid so far.
// The newly calculated balance is set on the loan.
func (l *Loan) BalanceRecalc() float64 {
	rem := l.balance

	for _, payment := range l.payments {
		rem -= payment
	}

	l.balance = max(0.0, rem)
	return l.balance
}

// Principal returns the original

// ApplyPayment records a payment. It rejects non-positive amounts and
// payments larger than the balance, leaving the loan unchanged.
func (l *Loan) ApplyPayment(amount float64) error {
	if amount <= 0 {
		return ErrInvalidPayment
	}
	if amount > l.balance {
		return ErrOverpayment
	}

	l.payments = append(l.payments, amount)
	l.balance -= amount
	return nil
}

// PercentPaid returns a fractional float64 indicating how much of the loan
// is paid off. Returns a value between 0.0 and 1.0 (inclusive).
func (l *Loan) PercentPaid() float64 {
	return min(1.0, max(0.0, (l.principal-l.balance)/l.principal))
}

// IsPaidOff reports whether the balance has reached zero.
func (l *Loan) IsPaidOff() bool {
	return l.balance == 0.0
}

// Payments returns the accepted payments in the order they were made.
// Changing the returned slice must not change the loan (return a copy/clone).
func (l *Loan) Payments() []float64 {
	return slices.Clone(l.payments)
}

func main() {
	l, err := NewLoan("L-100", 1200)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Initial balance: %.2f/%.2f, %% paid: %.1f\n", l.Balance(), l.Principal(), l.PercentPaid()*100)

	for _, amount := range []float64{200, 500, 800, 500} {
		if err := l.ApplyPayment(amount); err != nil {
			fmt.Printf("payment %.2f rejected: %v\n", amount, err)
			continue
		}
		fmt.Printf("paid %.2f, balance %.2f/%.2f, %% paid: %.1f\n", amount, l.Balance(), l.Principal(), l.PercentPaid()*100)
	}
	fmt.Printf("%s paid off: %v, payments: %v, %% paid: %.1f, final calc balance: %.2f\n", l.ID(), l.IsPaidOff(), l.Payments(), l.PercentPaid()*100, l.BalanceRecalc())
}
