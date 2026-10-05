package main

import (
	"errors"
	"strings"
	"testing"
)

const wrapPrefix = "invalid application: "

func TestValidateBorrowerMissing(t *testing.T) {
	if err := ValidateBorrower(""); !errors.Is(err, ErrMissingBorrower) {
		t.Errorf("ValidateBorrower(%q) = %v, want ErrMissingBorrower", "", err)
	}
	if err := ValidateBorrower("   "); !errors.Is(err, ErrMissingBorrower) {
		t.Errorf("ValidateBorrower(%q) = %v, want ErrMissingBorrower", "   ", err)
	}
	if err := ValidateBorrower("\t\n"); !errors.Is(err, ErrMissingBorrower) {
		t.Errorf("ValidateBorrower(%q) = %v, want ErrMissingBorrower", "\t\n", err)
	}
}

func TestValidateBorrowerPresent(t *testing.T) {
	if err := ValidateBorrower("Chris Bowe"); err != nil {
		t.Errorf("ValidateBorrower(%q) = %v, want nil", "Chris Bowe", err)
	}
}

func TestValidateAmountInRange(t *testing.T) {
	if err := ValidateAmount(1000); err != nil {
		t.Errorf("ValidateAmount(1000) = %v, want nil (Min is inclusive)", err)
	}
	if err := ValidateAmount(1000000); err != nil {
		t.Errorf("ValidateAmount(1000000) = %v, want nil (Max is inclusive)", err)
	}
	if err := ValidateAmount(250000); err != nil {
		t.Errorf("ValidateAmount(250000) = %v, want nil", err)
	}
}

func TestValidateAmountOutOfRange(t *testing.T) {
	err := ValidateAmount(999.99)
	if err == nil {
		t.Fatal("ValidateAmount(999.99) = nil, want a *RangeError")
	}

	var re *RangeError
	if !errors.As(err, &re) {
		t.Fatalf("ValidateAmount(999.99) = %v (%T), want a *RangeError", err, err)
	}
	if re.Field != "amount" {
		t.Errorf("Field = %q, want %q", re.Field, "amount")
	}
	if re.Value != 999.99 {
		t.Errorf("Value = %v, want 999.99", re.Value)
	}
	if re.Min != MinAmount || re.Max != MaxAmount {
		t.Errorf("Min, Max = %v, %v; want %v, %v", re.Min, re.Max, MinAmount, MaxAmount)
	}

	if err := ValidateAmount(1000000.01); err == nil {
		t.Error("ValidateAmount(1000000.01) = nil, want a *RangeError")
	}
	if err := ValidateAmount(0); err == nil {
		t.Error("ValidateAmount(0) = nil, want a *RangeError")
	}
}

func TestValidateTerm(t *testing.T) {
	if err := ValidateTerm(12); err != nil {
		t.Errorf("ValidateTerm(12) = %v, want nil", err)
	}
	if err := ValidateTerm(480); err != nil {
		t.Errorf("ValidateTerm(480) = %v, want nil", err)
	}

	err := ValidateTerm(6)
	if err == nil {
		t.Fatal("ValidateTerm(6) = nil, want a *RangeError")
	}

	var re *RangeError
	if !errors.As(err, &re) {
		t.Fatalf("ValidateTerm(6) = %v (%T), want a *RangeError", err, err)
	}
	if re.Field != "termMonths" {
		t.Errorf("Field = %q, want %q", re.Field, "termMonths")
	}
	if re.Value != 6 {
		t.Errorf("Value = %v, want 6", re.Value)
	}

	if err := ValidateTerm(481); err == nil {
		t.Error("ValidateTerm(481) = nil, want a *RangeError")
	}
}

func TestValidateApplicationValid(t *testing.T) {
	if err := ValidateApplication("Chris Bowe", 250000, 360); err != nil {
		t.Errorf("ValidateApplication(valid input) = %v, want nil", err)
	}
}

func TestValidateApplicationWrapsSentinel(t *testing.T) {
	err := ValidateApplication("", 250000, 360)
	if err == nil {
		t.Fatal("ValidateApplication(no borrower) = nil, want an error")
	}
	if !errors.Is(err, ErrMissingBorrower) {
		t.Errorf("errors.Is(err, ErrMissingBorrower) = false for %v; wrap with %%w, not %%v", err)
	}
	if !strings.HasPrefix(err.Error(), wrapPrefix) {
		t.Errorf("err.Error() = %q, want it to start with %q", err.Error(), wrapPrefix)
	}
}

func TestValidateApplicationWrapsRangeError(t *testing.T) {
	err := ValidateApplication("Chris Bowe", 50, 360)
	if err == nil {
		t.Fatal("ValidateApplication(amount 50) = nil, want an error")
	}

	var re *RangeError
	if !errors.As(err, &re) {
		t.Fatalf("errors.As(err, &re) = false for %v (%T); wrap with %%w, not %%v", err, err)
	}
	if re.Field != "amount" {
		t.Errorf("Field = %q, want %q", re.Field, "amount")
	}
	if !strings.HasPrefix(err.Error(), wrapPrefix) {
		t.Errorf("err.Error() = %q, want it to start with %q", err.Error(), wrapPrefix)
	}
}

// Checks run in a fixed order and stop at the first failure.
func TestValidateApplicationChecksInOrder(t *testing.T) {
	// Everything is wrong: the borrower error comes back.
	if err := ValidateApplication("", 50, 1); !errors.Is(err, ErrMissingBorrower) {
		t.Errorf("all-invalid input returned %v, want ErrMissingBorrower first", err)
	}

	// Borrower is fine, amount and term are not: the amount error comes back.
	err := ValidateApplication("Chris Bowe", 50, 1)
	var re *RangeError
	if !errors.As(err, &re) {
		t.Fatalf("ValidateApplication(bad amount and term) = %v, want a *RangeError", err)
	}
	if re.Field != "amount" {
		t.Errorf("Field = %q, want %q (amount is checked before termMonths)", re.Field, "amount")
	}
}

func TestRangeErrorMessage(t *testing.T) {
	err := ValidateAmount(50)
	if err == nil {
		t.Fatal("ValidateAmount(50) = nil, want a *RangeError")
	}
	msg := err.Error()
	if !strings.Contains(msg, "amount") {
		t.Errorf("message %q does not name the field", msg)
	}
	if !strings.Contains(msg, "50") {
		t.Errorf("message %q does not include the value", msg)
	}
}
