package main

import (
	"math"
	"testing"
)

// closeEnough compares dollar amounts to the cent.
// fractional cents are considered "close enough"
func closeEnough(got, want float64) bool {
	return math.Abs(got-want) < 0.01
}

func TestMonthlyPayment30Year(t *testing.T) {
	payment, interest := MonthlyPayment(250000, 0.065, 360)

	if !closeEnough(payment, 1580.17) {
		t.Errorf("payment = %.4f, want 1580.17", payment)
	}
	if !closeEnough(interest, 318861.22) {
		t.Errorf("totalInterest = %.4f, want 318861.22", interest)
	}
}

func TestMonthlyPayment5YearAuto(t *testing.T) {
	payment, interest := MonthlyPayment(30000, 0.0499, 60)

	if !closeEnough(payment, 566.00) {
		t.Errorf("payment = %.4f, want 566.00", payment)
	}
	if !closeEnough(interest, 3959.97) {
		t.Errorf("totalInterest = %.4f, want 3959.97", interest)
	}
}

func TestMonthlyPayment15Year(t *testing.T) {
	payment, interest := MonthlyPayment(500000, 0.0725, 180)

	if !closeEnough(payment, 4564.31) {
		t.Errorf("payment = %.4f, want 4564.31", payment)
	}
	if !closeEnough(interest, 321576.59) {
		t.Errorf("totalInterest = %.4f, want 321576.59", interest)
	}
}

// A zero rate must not divide by zero: it is a plain principal split.
func TestMonthlyPaymentZeroRate(t *testing.T) {
	payment, interest := MonthlyPayment(12000, 0, 24)

	if !closeEnough(payment, 500) {
		t.Errorf("payment = %.4f, want 500.00", payment)
	}
	if !closeEnough(interest, 0) {
		t.Errorf("totalInterest = %.4f, want 0.00", interest)
	}
}

func TestMonthlyPaymentGuards(t *testing.T) {
	if payment, interest := MonthlyPayment(250000, 0.065, 0); payment != 0 || interest != 0 {
		t.Errorf("termMonths 0: got (%v, %v), want (0, 0)", payment, interest)
	}
	if payment, interest := MonthlyPayment(250000, 0.065, -12); payment != 0 || interest != 0 {
		t.Errorf("negative term: got (%v, %v), want (0, 0)", payment, interest)
	}
	if payment, interest := MonthlyPayment(0, 0.065, 360); payment != 0 || interest != 0 {
		t.Errorf("zero principal: got (%v, %v), want (0, 0)", payment, interest)
	}
	if payment, interest := MonthlyPayment(-5000, 0.065, 360); payment != 0 || interest != 0 {
		t.Errorf("negative principal: got (%v, %v), want (0, 0)", payment, interest)
	}
}

// The guards apply at every rate, including the zero-rate branch.
func TestMonthlyPaymentGuardsZeroRate(t *testing.T) {
	if payment, interest := MonthlyPayment(12000, 0, 0); payment != 0 || interest != 0 {
		t.Errorf("zero rate, termMonths 0: got (%v, %v), want (0, 0)", payment, interest)
	}
	if payment, interest := MonthlyPayment(-5000, 0, 24); payment != 0 || interest != 0 {
		t.Errorf("zero rate, negative principal: got (%v, %v), want (0, 0)", payment, interest)
	}
}

func TestTotalFeesNoExtras(t *testing.T) {
	if got := TotalFees(1200); !closeEnough(got, 1200) {
		t.Errorf("TotalFees(1200) = %.4f, want 1200.00", got)
	}
}

func TestTotalFeesWithExtras(t *testing.T) {
	if got := TotalFees(1200, 450, 85.50); !closeEnough(got, 1735.50) {
		t.Errorf("TotalFees(1200, 450, 85.50) = %.4f, want 1735.50", got)
	}
	if got := TotalFees(0, 1, 2, 3, 4, 5); !closeEnough(got, 15) {
		t.Errorf("TotalFees(0, 1..5) = %.4f, want 15.00", got)
	}
}

// A slice can be spread into a variadic parameter with ...
func TestTotalFeesSpreadSlice(t *testing.T) {
	extras := []float64{300, 125.25, 19.75}
	if got := TotalFees(1000, extras...); !closeEnough(got, 1445) {
		t.Errorf("TotalFees(1000, extras...) = %.4f, want 1445.00", got)
	}
}
