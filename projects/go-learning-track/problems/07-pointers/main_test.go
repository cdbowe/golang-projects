package main

import "testing"

func TestApplyPaymentMutates(t *testing.T) {
	l := Loan{ID: "L1", Balance: 1000, Status: "active"}

	ApplyPayment(&l, 300)

	if l.Balance != 700 {
		t.Errorf("Balance = %v, want 700 (the change must reach the caller's loan)", l.Balance)
	}
	if l.Status != "active" {
		t.Errorf("Status = %q, want %q", l.Status, "active")
	}
}

func TestApplyPaymentPaysOff(t *testing.T) {
	l := Loan{ID: "L1", Balance: 250, Status: "active"}

	ApplyPayment(&l, 250)

	if l.Balance != 0 || l.Status != "paid" {
		t.Errorf("exact payoff: got (%v, %q), want (0, %q)", l.Balance, l.Status, "paid")
	}
}

func TestApplyPaymentOverpayFloorsAtZero(t *testing.T) {
	l := Loan{ID: "L1", Balance: 250, Status: "active"}

	ApplyPayment(&l, 400)

	if l.Balance != 0 || l.Status != "paid" {
		t.Errorf("overpayment: got (%v, %q), want (0, %q)", l.Balance, l.Status, "paid")
	}
}

func TestApplyPaymentIgnoresNonPositive(t *testing.T) {
	l := Loan{ID: "L1", Balance: 1000, Status: "active"}

	ApplyPayment(&l, 0)
	ApplyPayment(&l, -50)

	if l.Balance != 1000 || l.Status != "active" {
		t.Errorf("after 0 and -50: got (%v, %q), want (1000, %q)", l.Balance, l.Status, "active")
	}
}

// A nil pointer must be a no-op, not a panic.
func TestApplyPaymentNil(t *testing.T) {
	ApplyPayment(nil, 100)
}

func TestWithPaymentLeavesOriginal(t *testing.T) {
	original := Loan{ID: "L1", Balance: 1000, Status: "active"}

	updated := WithPayment(original, 300)

	if original.Balance != 1000 {
		t.Errorf("original.Balance = %v, want 1000 (WithPayment must not change its input)", original.Balance)
	}
	if updated.Balance != 700 {
		t.Errorf("updated.Balance = %v, want 700", updated.Balance)
	}
	if updated.ID != "L1" || updated.Status != "active" {
		t.Errorf("updated = %+v, want ID %q and Status %q carried over", updated, "L1", "active")
	}
}

func TestWithPaymentPaysOff(t *testing.T) {
	original := Loan{ID: "L1", Balance: 100, Status: "active"}

	updated := WithPayment(original, 150)

	if updated.Balance != 0 || updated.Status != "paid" {
		t.Errorf("updated = (%v, %q), want (0, %q)", updated.Balance, updated.Status, "paid")
	}
	if original.Status != "active" {
		t.Errorf("original.Status = %q, want %q", original.Status, "active")
	}
}

func TestFindReturnsSliceElement(t *testing.T) {
	loans := []Loan{
		{ID: "L1", Balance: 100, Status: "active"},
		{ID: "L2", Balance: 200, Status: "active"},
	}

	p := Find(loans, "L2")
	if p == nil {
		t.Fatal("Find(loans, \"L2\") = nil, want a pointer")
	}
	if p != &loans[1] {
		t.Error("Find returned a pointer to a copy, not to loans[1]; is it the address of the range variable?")
	}

	p.Balance = 0
	if loans[1].Balance != 0 {
		t.Errorf("after writing through the pointer, loans[1].Balance = %v, want 0", loans[1].Balance)
	}
}

func TestFindMissing(t *testing.T) {
	loans := []Loan{{ID: "L1", Balance: 100, Status: "active"}}

	if p := Find(loans, "L9"); p != nil {
		t.Errorf("Find(loans, \"L9\") = %+v, want nil", *p)
	}
	if p := Find(nil, "L1"); p != nil {
		t.Errorf("Find(nil, \"L1\") = %+v, want nil", *p)
	}
}
