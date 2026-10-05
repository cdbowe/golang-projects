package main

import "errors"

// Errors returned by Quote. Compare with errors.Is.
var (
	ErrUnsupportedTerm = errors.New("unsupported term")
	ErrIneligible      = errors.New("borrower is not eligible")
	ErrRateUnavailable = errors.New("base rate unavailable")
)

// RateSource supplies today's base rate for a loan term.
type RateSource interface {
	BaseRate(termMonths int) (float64, error)
}
