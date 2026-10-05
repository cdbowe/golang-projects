package main

import (
	"log"
	"net/http"
)

// NewServer returns the API's routes. See README.md for the exact status
// codes, headers and error messages.
//
//	GET  /loans       list every loan
//	GET  /loans/{id}  one loan, or 404
//	POST /loans       create a loan from {"borrower_name": ..., "loan_amount": ...}
func NewServer(store *Store) http.Handler {
	mux := http.NewServeMux()
	// TODO: register the three routes (method and path in one pattern) and
	// write their handlers.
	return mux
}

func main() {
	store := NewStore()
	store.Create("Ana Ruiz", 250000)

	addr := "localhost:8080"
	log.Printf("listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, NewServer(store)))
}
