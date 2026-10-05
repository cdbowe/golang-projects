# go-learning-track

Go, from zero to a portfolio-grade API, for an engineer who already knows C#/.NET and TypeScript.
17 problems. Each one adds **exactly one** new concept, drives it with tests, and names the closest
C# analogue so existing instincts transfer instead of misfiring.

Spec of record: [`../LEARNER_SPEC.md`](../LEARNER_SPEC.md). Tutor rules: [`CLAUDE.md`](CLAUDE.md).

## How a problem works

| File | What it is |
|---|---|
| `README.md` | The task, with an acceptance-criteria table mapping each rule to a test name |
| `LESSON.md` | The teaching: why Go does it this way, syntax, **C# ↔ Go table**, worked example, gotchas |
| `HINTS.md` | Three progressive hints, collapsed. Open them in order |
| `main.go` | Starter. Compiles, returns zero values, carries `// TODO` |
| `main_test.go` | The contract. Fails on assertions until you solve it |

Read `LESSON.md`, then `README.md`, then write code. `HINTS.md` only after you're stuck.

Logic lives in exported functions; `main()` just calls them and prints. That holds from problem 01 on,
which is why every problem is testable without touching stdout.

## Run

```bash
go test ./problems/01-hello-world/...        # one problem
go run  ./problems/01-hello-world            # see it print
go test ./...                                # everything (expect reds ahead of where you are)
go test -v -run TestGreetingText ./...       # one test, verbose
gofmt -l . && go vet ./...                   # before you call anything done
go test -short ./...                         # skip tests that need Docker (15)
```

Go 1.27, module `github.com/cdbowe/go-learning-track`. Local only — the module path is an identifier,
nothing is fetched or pushed.

## Progress

### Fundamentals — 15-30 min each

- [ ] **01 hello-world** — `package main`, `fmt`, the `go run` / `go test` loop
- [ ] **02 variables-if** — `var` vs `:=`, zero values, `if`, `switch`
- [ ] **03 loops** — `for` as the only loop
- [ ] **04 functions** — multiple returns, named results, variadics
- [ ] **05 errors** — `error` values, `fmt.Errorf` + `%w`, `errors.Is`/`As`
- [ ] 06 slices-maps — slices, maps, comma-ok idiom
- [ ] 07 pointers — `&`, `*`, pass-by-value semantics
- [ ] 08 structs-methods — structs, value vs pointer receivers, `NewX`
- [ ] 09 packages-modules — sub-packages, exported vs unexported, `internal/`
- [ ] 10 interfaces — implicit interfaces, composition over inheritance

### Real-world Go — 45-90 min each

- [ ] 11 table-tests — table-driven tests, `t.Run`, `t.Helper`, `-cover`
- [ ] 12 json-generics — `encoding/json`, struct tags, generics
- [ ] 13 http-api — `net/http`, Go 1.22+ routing, `httptest`
- [ ] 14 concurrency — goroutines, channels, `sync.WaitGroup`, `context`
- [ ] 15 postgres — `pgx`, SQL migrations, Testcontainers-go
- [ ] 16 graphql — `gqlgen` schema-first, resolvers

### Capstone — may spill past a day

- [ ] 17 **loan task rules API** — a loan update arrives, rules decide which tasks to create, Postgres
  persists them in a transaction, REST and GraphQL share one service layer

Problems 01-15 are generated. 16 and the capstone come next. Nothing is padded to
fill time — if a problem takes 12 minutes, it was a 12-minute problem.

## What transfers from C#, and what doesn't

| You already know | In Go |
|---|---|
| xUnit + TDD | `testing` + `go test`, no assertion library — you write the `if` |
| `Task`/`async` | goroutines + channels, no `async` keyword, no colouring |
| `IEnumerable`/LINQ | `for` loops; there is no lazy query pipeline in stdlib |
| Exceptions | `error` as a returned value, checked with `if err != nil` |
| Interfaces + DI | implicit interfaces, defined by the *consumer*, usually one method |
| Inheritance | embedding and composition. There is no inheritance |
| `.csproj` + NuGet | one `go.mod` per module |
| Testcontainers for .NET | Testcontainers-go, same idea, arrives at problem 15 |
