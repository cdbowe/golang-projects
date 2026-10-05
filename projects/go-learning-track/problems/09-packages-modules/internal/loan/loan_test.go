// Black-box test: package loan_test is a separate package, so it sees only
// the exported API — exactly what any other importer sees.
package loan_test

import (
	"errors"
	"math"
	"slices"
	"testing"

	"github.com/cdbowe/go-learning-track/problems/09-packages-modules/internal/loan"
)

func closeEnough(got, want float64) bool {
	return math.Abs(got-want) < 0.01
}

func TestNew(t *testing.T) {
	l, err := loan.New("L1", 1000)
	if err != nil {
		t.Fatalf("loan.New returned error %v, want nil", err)
	}
	if l == nil {
		t.Fatal("loan.New returned a nil *Loan")
	}
	if l.ID() != "L1" {
		t.Errorf("ID() = %q, want %q", l.ID(), "L1")
	}
	if !closeEnough(l.Principal(), 1000) {
		t.Errorf("Principal() = %v, want 1000", l.Principal())
	}
	if !closeEnough(l.Balance(), 1000) {
		t.Errorf("Balance() = %v, want 1000", l.Balance())
	}
}

func TestNewErrors(t *testing.T) {
	if l, err := loan.New(" ", 1000); !errors.Is(err, loan.ErrMissingID) || l != nil {
		t.Errorf("loan.New(blank id) = (%v, %v), want (nil, ErrMissingID)", l, err)
	}
	if l, err := loan.New("L1", 0); !errors.Is(err, loan.ErrInvalidPrincipal) || l != nil {
		t.Errorf("loan.New(\"L1\", 0) = (%v, %v), want (nil, ErrInvalidPrincipal)", l, err)
	}
}

func TestApplyPayment(t *testing.T) {
	l, _ := loan.New("L1", 1000)
	if l == nil {
		t.Fatal("loan.New returned a nil *Loan")
	}

	if err := l.ApplyPayment(400); err != nil {
		t.Fatalf("ApplyPayment(400) error = %v, want nil", err)
	}
	if err := l.ApplyPayment(0); !errors.Is(err, loan.ErrInvalidPayment) {
		t.Errorf("ApplyPayment(0) error = %v, want ErrInvalidPayment", err)
	}
	if err := l.ApplyPayment(601); !errors.Is(err, loan.ErrOverpayment) {
		t.Errorf("ApplyPayment(601) error = %v, want ErrOverpayment", err)
	}
	if !closeEnough(l.Balance(), 600) {
		t.Errorf("Balance() = %v, want 600", l.Balance())
	}

	if err := l.ApplyPayment(600); err != nil {
		t.Fatalf("ApplyPayment(600) error = %v, want nil", err)
	}
	if !l.IsPaidOff() {
		t.Error("IsPaidOff() = false after paying the full balance, want true")
	}
}

func TestPayments(t *testing.T) {
	l, _ := loan.New("L1", 1000)
	if l == nil {
		t.Fatal("loan.New returned a nil *Loan")
	}
	_ = l.ApplyPayment(100)
	_ = l.ApplyPayment(250)

	got := l.Payments()
	if want := []float64{100, 250}; !slices.Equal(got, want) {
		t.Fatalf("Payments() = %v, want %v", got, want)
	}

	got[0] = 999
	if l.Payments()[0] != 100 {
		t.Error("editing the result of Payments() changed the loan; return a copy")
	}
}
