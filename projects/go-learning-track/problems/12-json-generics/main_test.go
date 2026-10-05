package main

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// assertLoan compares every JSON-mapped field. InternalNotes is not compared.
func assertLoan(t *testing.T, got, want Loan) {
	t.Helper()
	if got.ID != want.ID || got.BorrowerName != want.BorrowerName ||
		got.LoanAmount != want.LoanAmount || got.Status != want.Status ||
		!slices.Equal(got.Documents, want.Documents) {
		t.Errorf("loan mismatch:\n got  %+v\n want %+v", got, want)
	}
}

func TestParseLoansFile(t *testing.T) {
	f, err := os.Open("testdata/loans.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })

	loans, err := ParseLoans(f)
	if err != nil {
		t.Fatalf("ParseLoans error = %v, want nil", err)
	}
	if len(loans) != 3 {
		t.Fatalf("got %d loans, want 3 (are the struct tags right?)", len(loans))
	}

	want := []Loan{
		{ID: "L-100", BorrowerName: "Ana Ruiz", LoanAmount: 250000, Status: "open", Documents: []string{"paystub", "w2"}},
		{ID: "L-101", BorrowerName: "Ben Ode", LoanAmount: 480000.5, Status: "closing"},
		{ID: "L-102", BorrowerName: "Cy Park", LoanAmount: 95000, Status: "open"},
	}
	for i := range want {
		assertLoan(t, loans[i], want[i])
	}
}

func TestParseLoansEmptyArray(t *testing.T) {
	loans, err := ParseLoans(strings.NewReader(`[]`))
	if err != nil {
		t.Fatalf("ParseLoans(`[]`) error = %v, want nil", err)
	}
	if len(loans) != 0 {
		t.Errorf("got %d loans, want 0", len(loans))
	}
}

func TestParseLoansErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		typeErr bool // expect a *json.UnmarshalTypeError somewhere in the chain
	}{
		{"truncated", `[{"id": "L-1"`, false},
		{"unknown key", `[{"id": "L-1", "nickname": "x"}]`, false},
		{"amount is a string", `[{"id": "L-1", "loan_amount": "lots"}]`, true},
		{"object, not array", `{"id": "L-1"}`, true},
		{"empty input", ``, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseLoans(strings.NewReader(tc.input))
			if err == nil {
				t.Fatalf("ParseLoans(%q) error = nil, want an error", tc.input)
			}
			if !strings.HasPrefix(err.Error(), "parse loans: ") {
				t.Errorf("error %q does not start with %q", err.Error(), "parse loans: ")
			}
			var typeErr *json.UnmarshalTypeError
			if tc.typeErr && !errors.As(err, &typeErr) {
				t.Errorf("error %v: want a *json.UnmarshalTypeError in the chain (wrap with %%w)", err)
			}
		})
	}
}

func TestEncodeLoan(t *testing.T) {
	tests := []struct {
		name string
		loan Loan
		want string
	}{
		{
			name: "no documents",
			loan: Loan{ID: "L-1", BorrowerName: "Ana", LoanAmount: 250000, Status: "open"},
			want: `{"id":"L-1","borrower_name":"Ana","loan_amount":250000,"status":"open"}`,
		},
		{
			name: "empty, non-nil documents are omitted too",
			loan: Loan{ID: "L-1", BorrowerName: "Ana", LoanAmount: 250000, Status: "open", Documents: []string{}},
			want: `{"id":"L-1","borrower_name":"Ana","loan_amount":250000,"status":"open"}`,
		},
		{
			name: "with documents",
			loan: Loan{ID: "L-2", BorrowerName: "Ben", LoanAmount: 99.5, Status: "closing", Documents: []string{"w2"}},
			want: `{"id":"L-2","borrower_name":"Ben","loan_amount":99.5,"status":"closing","documents":["w2"]}`,
		},
		{
			name: "internal notes never leave",
			loan: Loan{ID: "L-3", BorrowerName: "Cy", LoanAmount: 1, Status: "open", InternalNotes: "VIP, do not share"},
			want: `{"id":"L-3","borrower_name":"Cy","loan_amount":1,"status":"open"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := EncodeLoan(tc.loan)
			if err != nil {
				t.Fatalf("EncodeLoan error = %v, want nil", err)
			}
			if string(got) != tc.want {
				t.Errorf("EncodeLoan =\n %s\nwant\n %s", got, tc.want)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	original := Loan{ID: "L-9", BorrowerName: "Dee", LoanAmount: 123456.78, Status: "open", Documents: []string{"w2", "appraisal"}, InternalNotes: "secret"}

	encoded, err := EncodeLoan(original)
	if err != nil {
		t.Fatalf("EncodeLoan error = %v", err)
	}
	decoded, err := ParseLoans(strings.NewReader("[" + string(encoded) + "]"))
	if err != nil {
		t.Fatalf("ParseLoans error = %v (input %s)", err, encoded)
	}
	if len(decoded) != 1 {
		t.Fatalf("got %d loans, want 1", len(decoded))
	}

	assertLoan(t, decoded[0], original)
	if decoded[0].InternalNotes != "" {
		t.Errorf("InternalNotes = %q after a round trip, want empty", decoded[0].InternalNotes)
	}
}

func TestFilter(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	original := slices.Clone(nums)

	even := Filter(nums, func(n int) bool { return n%2 == 0 })
	if want := []int{2, 4, 6}; !slices.Equal(even, want) {
		t.Errorf("Filter(even) = %v, want %v", even, want)
	}
	if !slices.Equal(nums, original) {
		t.Errorf("input changed to %v, want %v (Filter must not reuse the input's backing array)", nums, original)
	}

	none := Filter(nums, func(n int) bool { return n > 100 })
	if len(none) != 0 {
		t.Errorf("Filter(no matches) = %v, want empty", none)
	}

	loans := []Loan{{ID: "A", Status: "open"}, {ID: "B", Status: "closing"}, {ID: "C", Status: "open"}}
	open := Filter(loans, func(l Loan) bool { return l.Status == "open" })
	if len(open) != 2 || open[0].ID != "A" || open[1].ID != "C" {
		t.Errorf("Filter(open loans) = %+v, want loans A and C in order", open)
	}
}

func TestMap(t *testing.T) {
	loans := []Loan{{ID: "A"}, {ID: "B"}}
	if got, want := Map(loans, func(l Loan) string { return l.ID }), []string{"A", "B"}; !slices.Equal(got, want) {
		t.Errorf("Map(loan IDs) = %v, want %v", got, want)
	}

	if got, want := Map([]int{1, 22, 333}, strconv.Itoa), []string{"1", "22", "333"}; !slices.Equal(got, want) {
		t.Errorf("Map(strconv.Itoa) = %v, want %v", got, want)
	}

	if got := Map([]int{}, strconv.Itoa); len(got) != 0 {
		t.Errorf("Map(empty) = %v, want empty", got)
	}
}

// cents is a named type whose underlying type is int.
type cents int

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Errorf("Sum(ints) = %v, want 6", got)
	}
	if got := Sum([]float64{0.5, 0.25}); got != 0.75 {
		t.Errorf("Sum(float64s) = %v, want 0.75", got)
	}
	if got := Sum([]cents{150, 250}); got != 400 {
		t.Errorf("Sum(cents) = %v, want 400", got)
	}
	if got := Sum([]int{}); got != 0 {
		t.Errorf("Sum(empty) = %v, want 0", got)
	}
}
