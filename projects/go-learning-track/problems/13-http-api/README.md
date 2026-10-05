# 13 - Loans REST API

**New concept:** `net/http` — handlers, method + path routing (Go 1.22+ `ServeMux`), JSON responses, `httptest`
**Builds on:** 05, 10, 11, 12

## Task

`Store` (in-memory, given complete) and `Loan` are in `store.go`. Implement `NewServer` in `main.go`:
register three routes and write their handlers.

| Route | Success | Failure |
|---|---|---|
| `GET /loans` | `200`, JSON array of `store.List()`. Empty store → `[]` | — |
| `GET /loans/{id}` | `200`, the loan | `404` `{"error":"loan not found"}` |
| `POST /loans` | `201`, the created loan, header `Location: /loans/<id>` | `400` (below) |

`POST` body: `{"borrower_name": "Ana Ruiz", "loan_amount": 250000}`.

| Bad input | `400` body |
|---|---|
| malformed JSON, wrong types, unknown fields, empty body | `{"error":"invalid JSON body"}` |
| `borrower_name` missing, empty or whitespace | `{"error":"borrower_name is required"}` |
| `loan_amount` zero or negative | `{"error":"loan_amount must be positive"}` |

Every JSON response sets `Content-Type: application/json` exactly. Other methods on these paths
(`DELETE /loans`, `PUT /loans/{id}`) return `405` — you get this free if you register routes the 1.22
way.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | Empty list is `[]`, not `null` | `TestListEmpty` |
| 2 | List returns every loan in ID order | `TestListLoans` |
| 3 | Get by ID; unknown ID is a JSON 404 | `TestGetLoan`, `TestGetLoanNotFound` |
| 4 | Create returns 201, the loan, `Location`, and persists it | `TestCreateLoan` |
| 5 | Each bad input gets its exact 400 message, and nothing is stored | `TestCreateLoanRejectsBadInput` |
| 6 | Wrong method → 405 with an `Allow` header | `TestMethodNotAllowed` |
| 7 | Works over a real TCP connection | `TestOverRealHTTP` |

## Run

```bash
go test -v ./problems/13-http-api/...
go run ./problems/13-http-api          # then, in another terminal:
curl -s localhost:8080/loans
curl -s -i -X POST localhost:8080/loans -d '{"borrower_name":"Ben","loan_amount":1200}'
curl -s -i localhost:8080/loans/L-999
```

## Stretch (optional)

Add `GET /healthz` returning `200` `{"status":"ok"}`, with a test in the existing style.
