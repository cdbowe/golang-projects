package main

import (
	"fmt"
	"strings"
)

// Label returns the servicing label for one month of a payment schedule.
// See the table in README.md. Months below 1 are "Invalid".
func Label(month int) string {
	switch {
	case month < 1:
		return "Invalid"
	case month%3 == 0 && month%4 == 0:
		return "Escrow check + Statement audit"
	case month%3 == 0:
		return "Escrow check"
	case month%4 == 0:
		return "Statement audit"
	default:
		return "Payment"
	}
}

// Schedule returns one line per month, "N: Label", joined with "\n" and with
// no trailing newline. months <= 0 returns the empty string.
func Schedule(months int) string {
	var sb strings.Builder

	for i := 1; i <= months; i++ {
		if i > 1 {
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, "%d: %s", i, Label(i)) // equivalent to: sb.WriteString(fmt.Sprintf("%v: %v", i, Label(i)))
	}

	return sb.String()
}

// Same as the Schedule method, but using += operator
func Schedule2(months int) string {
	var s string

	for i := 1; i <= months; i++ {
		if i > 1 {
			s += "\n"
		}
		s += fmt.Sprintf("%d: %s", i, Label(i))
	}

	return s
}

func main() {
	fmt.Println(Schedule(12))
}
