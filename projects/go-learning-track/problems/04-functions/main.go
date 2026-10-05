package main

import (
	"fmt"
	"math"
)

// MonthlyPayment returns the level monthly payment for an amortizing loan and
// the total interest paid over the full term.
//
// principal is in dollars, annualRate is a decimal fraction (0.065 == 6.5%),
// termMonths is the number of payments. See README.md for the formula and the
// guard conditions.
func MonthlyPayment(principal, annualRate float64, termMonths int) (payment, totalInterest float64) {
	// Term months and principal must both be positive
	if principal <= 0.0 || termMonths <= 0 {
		return 0, 0
	}

	// No annual rate means the calculation is simple
	if annualRate == 0.0 {
		return principal / float64(termMonths), 0
	}

	r := annualRate / 12

	// payment  = principal * r / (1 - (1+r)^-termMonths)
	payment = principal * r / (1 - math.Pow(1+r, -float64(termMonths)))
	totalInterest = payment*float64(termMonths) - principal

	// bare returns reads slightly better to me because it signals the presence of named result variables
	// equivalent to: return payment, totalInterest
	return
}

// TotalFees returns base plus every extra fee passed in. With no extras it
// returns base.
func TotalFees(base float64, extra ...float64) float64 {
	fee := base

	for _, v := range extra {
		fee += v
	}

	return fee
}

func main() {
	principal := 250000.00
	annualRate := 0.065
	termMonths := 360

	fmt.Printf("principal:  $%.2f | annual rate:  %.2f%% | term: %v months\n", principal, (annualRate * 100), termMonths)

	payment, interest := MonthlyPayment(principal, annualRate, termMonths)
	baseFee := 1200.00
	extra := []float64{450, 85.50}
	fees := TotalFees(baseFee, extra...)

	fmt.Printf("payment:  $%.2f\n", payment)
	fmt.Printf("interest: $%.2f\n", interest)
	fmt.Printf("fees:     $%.2f\n", fees)
}
