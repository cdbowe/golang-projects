package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Loan is one record in a loans JSON file. The Go field names differ from the
// JSON keys on purpose: struct tags connect them. See README.md for the keys.
type Loan struct {
	// TODO: add a struct tag to every field.
	ID            string
	BorrowerName  string
	LoanAmount    float64
	Status        string
	Documents     []string
	InternalNotes string // never read from or written to JSON
}

// ParseLoans decodes a JSON array of loans from r. Unknown keys are an error.
// Every error is wrapped with the prefix "parse loans: ".
func ParseLoans(r io.Reader) ([]Loan, error) {
	// TODO
	return nil, nil
}

// EncodeLoan returns l as compact JSON. Documents is left out when empty, and
// InternalNotes is never included.
func EncodeLoan(l Loan) ([]byte, error) {
	// TODO
	return nil, nil
}

// Number is the constraint for Sum.
type Number interface {
	// TODO: list the types Sum accepts: ints and float64s, including named
	// types built on them (like `type Cents int`).
}

// Filter returns the items for which keep returns true, in order. It never
// modifies items.
func Filter[T any](items []T, keep func(T) bool) []T {
	// TODO
	return nil
}

// Map returns f applied to every item, in order.
func Map[T, U any](items []T, f func(T) U) []U {
	// TODO
	return nil
}

// Sum adds up nums. An empty slice sums to zero.
func Sum[N Number](nums []N) N {
	var total N
	// TODO: add every number to total. (This line won't compile until Number
	// lists types that support +.)
	return total
}

func main() {
	path := "problems/12-json-generics/testdata/loans.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	loans, err := ParseLoans(bytes.NewReader(data))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	open := Filter(loans, func(l Loan) bool { return l.Status == "open" })
	fmt.Println("open loans:", Map(open, func(l Loan) string { return l.ID }))
	fmt.Printf("total amount: $%.2f\n", Sum(Map(loans, func(l Loan) float64 { return l.LoanAmount })))

	if len(loans) > 0 {
		out, err := EncodeLoan(loans[0])
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println("first loan as JSON:", string(out))
	}
}
