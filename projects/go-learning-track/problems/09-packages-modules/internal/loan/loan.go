// Package loan models a loan and the payments made against it.
package loan

import "slices"

// Loan tracks a principal and the payments made against it. Create one with
// New; the zero value is not usable.
type Loan struct {
	id        string
	principal float64
	balance   float64
	payments  []float64
}

// New validates its inputs and returns a ready-to-use loan. A blank id returns
// ErrMissingID; a non-positive principal returns ErrInvalidPrincipal.
func New(id string, principal float64) (*Loan, error) {
	if !validID(id) {
		return nil, ErrMissingID
	}
	if !validAmount(principal) {
		return nil, ErrInvalidPrincipal
	}

	return &Loan{id: id, principal: principal, balance: principal}, nil
}

// ID returns the loan's identifier.
func (l *Loan) ID() string {
	return l.id
}

// Principal returns the original amount borrowed.
func (l *Loan) Principal() float64 {
	return l.principal
}

// Balance returns the current outstanding balance.
func (l *Loan) Balance() float64 {
	return max(0.0, l.balance)
}

// BalanceRecalc returns the principal minus everything paid so far.
func (l *Loan) BalanceRecalc() float64 {
	rem := l.principal
	for _, payment := range l.payments {
		rem -= payment
	}
	return max(0.0, rem)
}

// ApplyPayment records a payment. A non-positive amount returns
// ErrInvalidPayment; more than the balance returns ErrOverpayment. Rejected
// payments leave the loan unchanged.
func (l *Loan) ApplyPayment(amount float64) error {
	if !validAmount(amount) {
		return ErrInvalidPayment
	}
	if amount > l.balance {
		return ErrOverpayment
	}

	l.payments = append(l.payments, amount)
	l.balance -= amount
	return nil
}

// IsPaidOff reports whether the balance has reached zero.
func (l *Loan) IsPaidOff() bool {
	return l.balance == 0.0
}

// Payments returns a copy of the accepted payments, oldest first.
func (l *Loan) Payments() []float64 {
	return slices.Clone(l.payments)
}
