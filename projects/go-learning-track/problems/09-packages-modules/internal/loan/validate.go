package loan

import (
	"strings"
)

// validID reports whether id is usable: not empty and not only whitespace.
// Unexported: other packages can't call it, and it isn't part of the API.
func validID(id string) bool {
	return strings.TrimSpace(id) != ""
}

// validAmount reports whether v is a usable money amount: strictly positive.
func validAmount(v float64) bool {
	return v > 0.0
}
