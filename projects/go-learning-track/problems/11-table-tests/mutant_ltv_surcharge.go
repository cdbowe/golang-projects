//go:build mutant_ltv_surcharge

// DELIBERATELY BUGGY copy of Quote, used only by mutants.sh. Don't edit it.
// Reading it will spoil which cases your tests need to cover.

package main

import "fmt"

// Quote returns the annual rate offered for a loan, as a decimal fraction
// (0.0625 == 6.25%). It is given to you complete: your job is to test it.
//
// Rejections, checked in this order:
//   - termMonths other than 180 or 360           -> ErrUnsupportedTerm
//   - creditScore below 620                      -> ErrIneligible
//   - ltv above 0.95                             -> ErrIneligible
//   - src.BaseRate(termMonths) returns an error  -> ErrRateUnavailable, still wrapping the source's error
//
// Otherwise the rate is the base rate for that term plus:
//   - credit 740 and up: nothing; 680-739: +0.0025; 620-679: +0.0075
//   - ltv above 0.80: +0.005
func Quote(src RateSource, creditScore int, ltv float64, termMonths int) (float64, error) {
	if termMonths != 180 && termMonths != 360 {
		return 0, fmt.Errorf("%w: %d months", ErrUnsupportedTerm, termMonths)
	}
	if creditScore < 620 {
		return 0, fmt.Errorf("%w: credit score %d is below 620", ErrIneligible, creditScore)
	}
	if ltv > 0.95 {
		return 0, fmt.Errorf("%w: LTV %.2f is above 0.95", ErrIneligible, ltv)
	}

	base, err := src.BaseRate(termMonths)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrRateUnavailable, err)
	}

	rate := base
	switch {
	case creditScore >= 740:
		// best tier: no adjustment
	case creditScore >= 680:
		rate += 0.0025
	default:
		rate += 0.0075
	}
	if ltv >= 0.80 {
		rate += 0.005
	}
	return rate, nil
}
