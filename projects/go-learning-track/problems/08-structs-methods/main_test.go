package main

import (
	"errors"
	"math"
	"slices"
	"testing"
)

func closeEnough(got, want float64) bool {
	return math.Abs(got-want) < 0.01
}

func TestNewLoan(t *testing.T) {
	l, err := NewLoan("L1", 1000)
	if err != nil {
		t.Fatalf("NewLoan returned error %v, want nil", err)
	}
	if l == nil {
		t.Fatal("NewLoan returned a nil *Loan")
	}
	if l.ID() != "L1" {
		t.Errorf("ID() = %q, want %q", l.ID(), "L1")
	}
	if !closeEnough(l.Balance(), 1000) {
		t.Errorf("Balance() = %v, want 1000", l.Balance())
	}
	if l.IsPaidOff() {
		t.Error("IsPaidOff() = true for a new loan, want false")
	}
	if len(l.Payments()) != 0 {
		t.Errorf("Payments() = %v, want none", l.Payments())
	}
}

func TestNewLoanMissingID(t *testing.T) {
	l, err := NewLoan("", 1000)
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("NewLoan(\"\", 1000) error = %v, want ErrMissingID", err)
	}
	if l != nil {
		t.Error("NewLoan returned a loan alongside an error, want nil")
	}

	if _, err := NewLoan("   ", 1000); !errors.Is(err, ErrMissingID) {
		t.Errorf("NewLoan(whitespace id) error = %v, want ErrMissingID", err)
	}
}

func TestNewLoanInvalidPrincipal(t *testing.T) {
	l, err := NewLoan("L1", 0)
	if !errors.Is(err, ErrInvalidPrincipal) {
		t.Errorf("NewLoan(\"L1\", 0) error = %v, want ErrInvalidPrincipal", err)
	}
	if l != nil {
		t.Error("NewLoan returned a loan alongside an error, want nil")
	}

	if _, err := NewLoan("L1", -500); !errors.Is(err, ErrInvalidPrincipal) {
		t.Errorf("NewLoan(\"L1\", -500) error = %v, want ErrInvalidPrincipal", err)
	}
}

// If this fails with the balance unchanged, look at ApplyPayment's receiver.
func TestApplyPaymentPersists(t *testing.T) {
	l, _ := NewLoan("L1", 1000)
	if l == nil {
		t.Fatal("NewLoan returned a nil *Loan")
	}

	if err := l.ApplyPayment(400); err != nil {
		t.Fatalf("ApplyPayment(400) error = %v, want nil", err)
	}
	if !closeEnough(l.Balance(), 600) {
		t.Errorf("Balance() after paying 400 = %v, want 600", l.Balance())
	}
}

func TestApplyPaymentRejectsNonPositive(t *testing.T) {
	l, _ := NewLoan("L1", 1000)
	if l == nil {
		t.Fatal("NewLoan returned a nil *Loan")
	}

	if err := l.ApplyPayment(0); !errors.Is(err, ErrInvalidPayment) {
		t.Errorf("ApplyPayment(0) error = %v, want ErrInvalidPayment", err)
	}
	if err := l.ApplyPayment(-10); !errors.Is(err, ErrInvalidPayment) {
		t.Errorf("ApplyPayment(-10) error = %v, want ErrInvalidPayment", err)
	}
	if !closeEnough(l.Balance(), 1000) || len(l.Payments()) != 0 {
		t.Errorf("after rejected payments: Balance %v, Payments %v; want 1000 and none", l.Balance(), l.Payments())
	}
}

func TestApplyPaymentRejectsOverpayment(t *testing.T) {
	l, _ := NewLoan("L1", 1000)
	if l == nil {
		t.Fatal("NewLoan returned a nil *Loan")
	}

	if err := l.ApplyPayment(1000.01); !errors.Is(err, ErrOverpayment) {
		t.Errorf("ApplyPayment(1000.01) error = %v, want ErrOverpayment", err)
	}
	if !closeEnough(l.Balance(), 1000) || len(l.Payments()) != 0 {
		t.Errorf("after rejected overpayment: Balance %v, Payments %v; want 1000 and none", l.Balance(), l.Payments())
	}
}

func TestPaidOff(t *testing.T) {
	l, _ := NewLoan("L1", 1000)
	if l == nil {
		t.Fatal("NewLoan returned a nil *Loan")
	}

	_ = l.ApplyPayment(400)
	if l.IsPaidOff() {
		t.Error("IsPaidOff() = true with 600 left, want false")
	}

	if err := l.ApplyPayment(600); err != nil {
		t.Fatalf("ApplyPayment(600) error = %v, want nil (paying the exact balance is allowed)", err)
	}
	if !l.IsPaidOff() {
		t.Error("IsPaidOff() = false after paying the full balance, want true")
	}
	if !closeEnough(l.Balance(), 0) {
		t.Errorf("Balance() = %v, want 0", l.Balance())
	}
}

func TestPaymentsHistory(t *testing.T) {
	l, _ := NewLoan("L1", 1000)
	if l == nil {
		t.Fatal("NewLoan returned a nil *Loan")
	}

	_ = l.ApplyPayment(100)
	_ = l.ApplyPayment(5000) // rejected: not recorded
	_ = l.ApplyPayment(250)

	if got, want := l.Payments(), []float64{100, 250}; !slices.Equal(got, want) {
		t.Errorf("Payments() = %v, want %v", got, want)
	}
}

func TestPaymentsReturnsCopy(t *testing.T) {
	l, _ := NewLoan("L1", 1000)
	if l == nil {
		t.Fatal("NewLoan returned a nil *Loan")
	}
	_ = l.ApplyPayment(100)

	p := l.Payments()
	if len(p) != 1 {
		t.Fatalf("Payments() = %v, want one payment", p)
	}
	p[0] = 999

	if got := l.Payments(); got[0] != 100 {
		t.Errorf("after editing the returned slice, Payments()[0] = %v, want 100 (return a copy)", got[0])
	}
}

// Two loans must not share state.
func TestLoansAreIndependent(t *testing.T) {
	a, _ := NewLoan("A", 1000)
	b, _ := NewLoan("B", 1000)
	if a == nil || b == nil {
		t.Fatal("NewLoan returned a nil *Loan")
	}

	_ = a.ApplyPayment(300)

	if !closeEnough(b.Balance(), 1000) {
		t.Errorf("b.Balance() = %v after paying a, want 1000", b.Balance())
	}
}
