package loan

import "errors"

// Errors returned by this package. Callers compare with errors.Is.
var (
	ErrMissingID        = errors.New("loan id is required")
	ErrInvalidPrincipal = errors.New("principal must be positive")
	ErrInvalidPayment   = errors.New("payment must be positive")
	ErrOverpayment      = errors.New("payment exceeds balance")
)
