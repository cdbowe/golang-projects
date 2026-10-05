package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// do sends one request straight to the handler, with no network involved.
func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, r)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// decode asserts a JSON response with the given status and decodes its body.
func decode[T any](t *testing.T, rec *httptest.ResponseRecorder, wantStatus int) T {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("status = %d, want %d (body: %q)", rec.Code, wantStatus, rec.Body.String())
	}
	// Result() is what the client receives: headers set after WriteHeader are lost.
	if ct := rec.Result().Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("body %q is not valid JSON for %T: %v", rec.Body.String(), v, err)
	}
	return v
}

// assertError asserts a JSON error response: {"error": "<msg>"}.
func assertError(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantMsg string) {
	t.Helper()
	body := decode[map[string]string](t, rec, wantStatus)
	if body["error"] != wantMsg {
		t.Errorf(`body["error"] = %q, want %q`, body["error"], wantMsg)
	}
}

func TestListEmpty(t *testing.T) {
	rec := do(t, NewServer(NewStore()), "GET", "/loans", "")

	loans := decode[[]Loan](t, rec, http.StatusOK)
	if len(loans) != 0 {
		t.Errorf("got %d loans, want 0", len(loans))
	}
	if got := strings.TrimSpace(rec.Body.String()); got != "[]" {
		t.Errorf("body = %q, want %q (an empty list, not null)", got, "[]")
	}
}

func TestListLoans(t *testing.T) {
	store := NewStore()
	store.Create("Ana", 100000)
	store.Create("Ben", 200000)

	loans := decode[[]Loan](t, do(t, NewServer(store), "GET", "/loans", ""), http.StatusOK)

	if len(loans) != 2 || loans[0].ID != "L-100" || loans[1].ID != "L-101" {
		t.Errorf("loans = %+v, want L-100 then L-101", loans)
	}
}

func TestGetLoan(t *testing.T) {
	store := NewStore()
	store.Create("Ana", 100000)

	got := decode[Loan](t, do(t, NewServer(store), "GET", "/loans/L-100", ""), http.StatusOK)

	want := Loan{ID: "L-100", BorrowerName: "Ana", LoanAmount: 100000, Status: "open"}
	if got != want {
		t.Errorf("loan = %+v, want %+v", got, want)
	}
}

func TestGetLoanNotFound(t *testing.T) {
	rec := do(t, NewServer(NewStore()), "GET", "/loans/L-999", "")
	assertError(t, rec, http.StatusNotFound, "loan not found")
}

func TestCreateLoan(t *testing.T) {
	store := NewStore()
	rec := do(t, NewServer(store), "POST", "/loans", `{"borrower_name": "Ana Ruiz", "loan_amount": 250000}`)

	got := decode[Loan](t, rec, http.StatusCreated)

	want := Loan{ID: "L-100", BorrowerName: "Ana Ruiz", LoanAmount: 250000, Status: "open"}
	if got != want {
		t.Errorf("created = %+v, want %+v", got, want)
	}
	if loc := rec.Result().Header.Get("Location"); loc != "/loans/L-100" {
		t.Errorf("Location = %q, want %q", loc, "/loans/L-100")
	}
	if _, ok := store.Get("L-100"); !ok {
		t.Error("the loan was not saved in the store")
	}
}

func TestCreateLoanRejectsBadInput(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantMsg string
	}{
		{"blank borrower", `{"borrower_name": "", "loan_amount": 1000}`, "borrower_name is required"},
		{"whitespace borrower", `{"borrower_name": "   ", "loan_amount": 1000}`, "borrower_name is required"},
		{"missing borrower", `{"loan_amount": 1000}`, "borrower_name is required"},
		{"zero amount", `{"borrower_name": "Ana", "loan_amount": 0}`, "loan_amount must be positive"},
		{"negative amount", `{"borrower_name": "Ana", "loan_amount": -5}`, "loan_amount must be positive"},
		{"malformed JSON", `{"borrower_name": "Ana",`, "invalid JSON body"},
		{"wrong type", `{"borrower_name": "Ana", "loan_amount": "lots"}`, "invalid JSON body"},
		{"unknown field", `{"borrower_name": "Ana", "loan_amount": 1000, "rate": 0.06}`, "invalid JSON body"},
		{"empty body", ``, "invalid JSON body"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := NewStore()
			rec := do(t, NewServer(store), "POST", "/loans", tc.body)

			assertError(t, rec, http.StatusBadRequest, tc.wantMsg)
			if n := len(store.List()); n != 0 {
				t.Errorf("store has %d loans after a rejected request, want 0", n)
			}
		})
	}
}

// Registering routes with their method gets you 405 responses for free.
func TestMethodNotAllowed(t *testing.T) {
	tests := []struct {
		method, path string
	}{
		{"DELETE", "/loans"},
		{"PUT", "/loans/L-100"},
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			rec := do(t, NewServer(NewStore()), tc.method, tc.path, "")
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
			}
			if allow := rec.Result().Header.Get("Allow"); !strings.Contains(allow, "GET") {
				t.Errorf("Allow = %q, want it to list GET", allow)
			}
		})
	}
}

// One test over a real TCP connection, the way a client would see it.
func TestOverRealHTTP(t *testing.T) {
	srv := httptest.NewServer(NewServer(NewStore()))
	t.Cleanup(srv.Close)

	resp, err := http.Post(srv.URL+"/loans", "application/json",
		strings.NewReader(`{"borrower_name": "Ana", "loan_amount": 1000}`))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	resp, err = http.Get(srv.URL + resp.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { resp.Body.Close() })

	var got Loan
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decoding GET response: %v", err)
	}
	if resp.StatusCode != http.StatusOK || got.BorrowerName != "Ana" {
		t.Errorf("GET Location: status %d, loan %+v; want 200 and Ana's loan", resp.StatusCode, got)
	}
}
