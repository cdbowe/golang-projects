# Hints - 13

<details>
<summary>Hint 1 - direction</summary>

Three handlers and two small helpers. Write a `writeJSON(w, status, v)` helper first, and an error
helper on top of it — every response in this task goes through one of them. Then each handler is:
read input, validate, call the store, write one response, `return`.

Handlers need the store. Either make them methods on a small struct holding `*Store`, or closures
inside `NewServer` that capture `store`.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- `mux.HandleFunc("GET /loans/{id}", ...)` and `r.PathValue("id")`.
- `json.NewDecoder(r.Body)` + `DisallowUnknownFields()` + `Decode(&in)`. Any decode error → `invalid JSON body`.
  An empty body makes `Decode` return `io.EOF` — also an error, so it's already covered.
- Request shape: a struct with `BorrowerName string` and `LoanAmount float64`, tagged like `Loan`.
- `strings.TrimSpace` for the borrower check (05).
- 201 + Location: `w.Header().Set("Location", "/loans/"+l.ID)` **before** `writeJSON`.
- Status constants: `http.StatusOK`, `http.StatusCreated`, `http.StatusBadRequest`, `http.StatusNotFound`.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
writeJSON(w, status, v):  set Content-Type; WriteHeader(status); json.NewEncoder(w).Encode(v)
writeError(w, status, msg): writeJSON(w, status, map[string]string{"error": msg})

NewServer(store):
    mux := NewServeMux()
    mux.HandleFunc("GET /loans", func(w, r) { writeJSON(w, 200, store.List()) })

    mux.HandleFunc("GET /loans/{id}", func(w, r) {
        l, ok := store.Get(r.PathValue("id"))
        if !ok { writeError(w, 404, "loan not found"); return }
        writeJSON(w, 200, l)
    })

    mux.HandleFunc("POST /loans", func(w, r) {
        decode into `in` with unknown fields disallowed → on error: 400 "invalid JSON body"; return
        if TrimSpace(in.BorrowerName) == "": 400 "borrower_name is required"; return
        if in.LoanAmount <= 0:               400 "loan_amount must be positive"; return
        l := store.Create(in.BorrowerName, in.LoanAmount)
        set Location; writeJSON(w, 201, l)
    })
    return mux
```

</details>
