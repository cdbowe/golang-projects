package main

import (
	"testing"

	"github.com/cdbowe/go-learning-track/problems/09-packages-modules/internal/loan"
)

func TestSummary(t *testing.T) {
	l, err := loan.New("L-100", 1000)
	if err != nil || l == nil {
		t.Fatalf("loan.New = (%v, %v), want a loan; finish internal/loan first", l, err)
	}
	_ = l.ApplyPayment(150)
	_ = l.ApplyPayment(250)

	want := "L-100: $600.00 of $1000.00 remaining, payments: 2"
	if got := Summary(l); got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}

func TestSummaryNewLoan(t *testing.T) {
	l, _ := loan.New("L-7", 250.5)
	if l == nil {
		t.Fatal("loan.New returned a nil *Loan; finish internal/loan first")
	}

	want := "L-7: $250.50 of $250.50 remaining, payments: 0"
	if got := Summary(l); got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}

func TestSummaryNil(t *testing.T) {
	if got := Summary(nil); got != "no loan" {
		t.Errorf("Summary(nil) = %q, want %q", got, "no loan")
	}
}
