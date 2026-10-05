# Lesson 13 - HTTP servers with net/http

## Why Go does it this way

- **The standard library is the framework.** `net/http` is production-grade: no Kestrel/ASP.NET split, no hosting model. Many Go services ship with zero web dependencies.
- **One interface does it all.** Anything with `ServeHTTP(w http.ResponseWriter, r *http.Request)` is an `http.Handler`. A mux is a handler; middleware is a handler wrapping a handler.
- **Routing grew up in Go 1.22.** Patterns carry the method and wildcards: `"GET /loans/{id}"`. Third-party routers are now optional.
- **Tests don't need a network.** `httptest.NewRecorder` captures what a handler writes; `httptest.NewServer` gives a real localhost URL when you want one.

## Syntax

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /offers", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // headers BEFORE status
	w.WriteHeader(http.StatusOK)                       // optional for 200; once only
	json.NewEncoder(w).Encode(offers)                  // body last
})

mux.HandleFunc("GET /offers/{lender}", func(w http.ResponseWriter, r *http.Request) {
	lender := r.PathValue("lender") // the {lender} wildcard
	_ = lender
})

mux.HandleFunc("POST /offers", func(w http.ResponseWriter, r *http.Request) {
	var in struct {                 // anonymous struct: request shape, used once
		Lender string  `json:"lender"`
		Rate   float64 `json:"rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest) // plain-text helper
		return                                            // ALWAYS return after writing an error
	}
})

log.Fatal(http.ListenAndServe("localhost:8080", mux))
```

Handlers as methods, sharing state without globals:

```go
type api struct{ store *Store }

func (a *api) list(w http.ResponseWriter, r *http.Request) { /* a.store... */ }

mux.HandleFunc("GET /loans", a.list) // a method value is a func
```

`defer f()` runs `f` when the surrounding function returns — used in `store.go` to guarantee
`Unlock`. Same job as `finally`.

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `http.Handler` / `HandlerFunc` | minimal API endpoint / `RequestDelegate` | One interface; no DI container, no model binding, no filters |
| `"GET /loans/{id}"` | `app.MapGet("/loans/{id}", ...)` | Method is part of the pattern string. Unmatched method → automatic 405 |
| `r.PathValue("id")` | route parameter binding | Always a `string`; convert yourself |
| `json.NewDecoder(r.Body).Decode` | `[FromBody]` | Manual. Validation is yours, too |
| `w.Header().Set` → `WriteHeader` → `Write` | `Results.Json(...)` | **Order matters.** Headers set after the first write are silently ignored |
| `http.Error` | `Results.Problem` | Writes `text/plain`. For JSON errors, write your own |
| `httptest.NewRecorder` | `WebApplicationFactory` + `HttpClient` | No host at all — call `ServeHTTP` directly |
| `defer` | `finally` / `using` | Runs at function exit, last-in first-out. Arguments are evaluated immediately |
| handler goroutine per request | thread-pool request | Every handler may run concurrently with others. Shared state needs locking (14) |

## Worked example

A JSON helper and a handler for a different resource:

```go
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err) // too late to change the status
	}
}

type rateAPI struct{ rates map[string]float64 }

func (a *rateAPI) get(w http.ResponseWriter, r *http.Request) {
	term := r.PathValue("term")
	rate, ok := a.rates[term]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no rate for " + term})
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"rate": rate})
}

// registered as: mux.HandleFunc("GET /rates/{term}", api.get)
```

## Gotchas for C# developers

- Forgetting `return` after writing an error: the handler keeps going and writes a second response. The log says `superfluous response.WriteHeader call`.
- `w.Header().Set` after `WriteHeader` or `Write` does nothing.
- `json.Encoder.Encode` appends `\n`. Clients don't care; exact-string tests do.
- A nil slice encodes as `null`. Lists should encode as `[]` — `Store.List` already returns non-nil.
- `DisallowUnknownFields` is per-decoder; set it before `Decode`.
- `"GET /loans"` also matches `HEAD`. `"/loans"` with no method matches **every** method — then 405 is your job.
- `http.ListenAndServe` blocks forever; it only returns an error.

## Check yourself

1. What does the client receive if a handler calls `writeJSON(w, 400, ...)` and then, missing a `return`, `writeJSON(w, 201, ...)`?
2. Why register `"POST /loans"` rather than `"/loans"` with a `switch r.Method`?
3. How does `TestCreateLoan` call your handler without opening a port?

<details>
<summary>Answers</summary>

1. Status 400 (the first `WriteHeader` wins), a body with **both** JSON documents concatenated, and a
   `superfluous WriteHeader` log line. The client probably fails to parse it.
2. The mux then knows which methods each path allows: it returns 405 with an `Allow` header for you,
   and the handler can't be reached with the wrong method.
3. `httptest.NewRequest` builds an `*http.Request`, `httptest.NewRecorder` implements
   `http.ResponseWriter`, and the test calls `handler.ServeHTTP(rec, req)` — an ordinary method call.

</details>

## Go deeper

- [pkg.go.dev/net/http#ServeMux](https://pkg.go.dev/net/http#ServeMux) — the pattern syntax
- [Go blog: Routing Enhancements for Go 1.22](https://go.dev/blog/routing-enhancements)
- [pkg.go.dev/net/http/httptest](https://pkg.go.dev/net/http/httptest)
- [Effective Go: defer](https://go.dev/doc/effective_go#defer)
