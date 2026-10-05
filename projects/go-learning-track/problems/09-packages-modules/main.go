package main

import (
	"fmt"

	"github.com/cdbowe/go-learning-track/problems/09-packages-modules/internal/loan"
)

// Summary describes a loan in one line for display. A nil loan is "no loan".
// Format: "L-100: $600.00 of $1000.00 remaining, payments: 2".
func Summary(l *loan.Loan) string {
	if l == nil {
		return "no loan"
	}

	return fmt.Sprintf("%s: $%.2f of $%.2f remaining, payments: %d", l.ID(), l.Balance(), l.Principal(), len(l.Payments()))
}

func main() {
	l, err := loan.New("L-100", 1200)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if err := l.ApplyPayment(200); err != nil {
		fmt.Println("payment rejected:", err)
	}
	fmt.Println(Summary(l))
	fmt.Println(Summary(nil))
}
