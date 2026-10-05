package main

import (
	"errors"
	"fmt"
	"strings"
)

// ErrMissingBorrower is a sentinel error: callers compare against it with
// errors.Is. It is created once, at package level, not per call.
var ErrMissingBorrower = errors.New("borrower is required")

// RangeError reports a numeric field that fell outside its allowed range.
// Callers recover the details with errors.As.
//
// This type is given to you — structs and methods are problem 08. All you do
// with it here is return it and wrap it.
type RangeError struct {
	Field string
	Value float64
	Min   float64
	Max   float64
}

// Error makes *RangeError satisfy the error interface.
func (e *RangeError) Error() string {
	return fmt.Sprintf("%s %v is outside the allowed range %v..%v", e.Field, e.Value, e.Min, e.Max)
}

// Allowed ranges for a loan application.
const (
	MinAmount = 1000.0
	MaxAmount = 1000000.0
	MinTerm   = 12
	MaxTerm   = 480
)

// ValidateBorrower returns ErrMissingBorrower when name is empty or only
// whitespace, and nil otherwise.
func ValidateBorrower(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrMissingBorrower
	}

	return nil
}

// ValidateAmount returns a *RangeError when amount is outside
// MinAmount..MaxAmount (inclusive), and nil otherwise.
func ValidateAmount(amount float64) error {
	if amount > MaxAmount || amount < MinAmount {
		return &RangeError{Field: "amount", Value: amount, Min: MinAmount, Max: MaxAmount}
	}

	return nil
}

// ValidateTerm returns a *RangeError when termMonths is outside
// MinTerm..MaxTerm (inclusive), and nil otherwise.
func ValidateTerm(termMonths int) error {
	if termMonths < MinTerm || termMonths > MaxTerm {
		return &RangeError{Field: "termMonths", Value: float64(termMonths), Min: MinTerm, Max: MaxTerm}
	}

	return nil
}

// ValidateApplication checks borrower, then amount, then term, and returns the
// first failure wrapped with the prefix "invalid application: ". It returns nil
// when everything is valid.
func ValidateApplication(borrower string, amount float64, termMonths int) error {
	if err := ValidateBorrower(borrower); err != nil {
		return fmt.Errorf("invalid application: %w", err)
	}

	if err := ValidateAmount(amount); err != nil {
		return fmt.Errorf("invalid application: %w", err)
	}

	if err := ValidateTerm(termMonths); err != nil {
		return fmt.Errorf("invalid application: %w", err)
	}

	return nil
}

func main() {
	if err := ValidateApplication("", 250000, 360); err != nil {
		fmt.Println("rejected:", err)
	}
	if err := ValidateApplication("Chris Bowe", 250000, 360); err == nil {
		fmt.Println("accepted: Chris Bowe, $250000, 360 months")
	}

	errs := []error{
		ValidateApplication("", 250000, 360),
		ValidateApplication("Foobar", 2500000, 360),
		ValidateApplication("Foobar", 250000, 36000),
		ValidateApplication("Foobar", 250000, 360)}

	var appErr *RangeError
	for i, v := range errs {
		fmt.Printf("ERR %d: ", (i + 1))
		if errors.Is(v, ErrMissingBorrower) {
			fmt.Println("INVALID BORROWER NAME")
		} else if errors.As(v, &appErr) {
			fmt.Printf("RANGE ERROR: %q: %v, Accepted Range: %v-%v\n", appErr.Field, appErr.Value, appErr.Min, appErr.Max)
		} else {
			fmt.Println("No error")
		}
	}

	// Use errors.Is() and errors.As() here, not in validator functions

	// var firstErr error
	// err := ValidateBorrower(borrower)
	// if err != nil {
	// 	firstErr = err
	// 	if errors.Is(err, ErrMissingBorrower) {
	// 		fmt.Printf("INVALID BORROWER NAME: %q\n", borrower)
	// 	} else {
	// 		fmt.Println("UNKNOWN ERROR TYPE FOR BORROWER")
	// 	}
	// }

	// if firstErr == nil {
	// 	err = ValidateAmount(amount)
	// 	if err != nil {
	// 		firstErr = err

	// 		var amtErr *RangeError
	// 		if errors.As(err, &amtErr) {
	// 			fmt.Printf("AMT ERROR: %q: %v, Accepted Range: %v-%v", amtErr.Field, amtErr.Value, amtErr.Min, amtErr.Max)
	// 		} else {
	// 			fmt.Println("UNKNOWN AMOUNT ERROR TYPE")
	// 		}
	// 	}
	// }

	// if firstErr == nil {
	// 	err = ValidateTerm(termMonths)
	// 	if err != nil {
	// 		firstErr = err

	// 		var termErr *RangeError
	// 		if errors.As(err, &termErr) {
	// 			fmt.Printf("TERM ERROR: %q: %v, Accepted Range: %v-%v", termErr.Field, termErr.Value, termErr.Min, termErr.Max)
	// 		} else {
	// 			fmt.Println("UNKNOWN TERM ERROR TYPE")
	// 		}
	// 	}
	// }

	// if firstErr != nil {
	// 	return fmt.Errorf("invalid application: %w", firstErr)
	// }

	// return nil
}
