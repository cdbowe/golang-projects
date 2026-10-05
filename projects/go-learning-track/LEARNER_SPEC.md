# golang-projects: Problem-Set Spec

Give this file to Claude in the `golang-projects` repo. Claude generates the problem set, lessons, hints, starter code and tests described here. **Claude does not solve the problems.**

## Learner

| Fact | Implication |
|---|---|
| Chris Bowe, senior C#/.NET + TypeScript/React engineer, 14 yrs | Skip "what is a variable". Teach Go's *differences* from C#, fast |
| Zero Go experience | Start at Hello World anyway; steps are small, pace is fast |
| Strong xUnit/TDD habits, uses Testcontainers in .NET | Tests drive every problem; Testcontainers-go in the DB problem |
| ADHD/ASD | One new concept per problem. Short, scannable lessons. Explicit acceptance criteria. No walls of prose |
| Goal | Weekend sprint ending in a portfolio-grade capstone: a loan task rules API (Go + Postgres + GraphQL) |

## Repo layout

```
golang-projects/
├── go.mod                      # module github.com/cdbowe/golang-projects
├── README.md                   # track overview, progress checklist, how to run
├── CLAUDE.md                   # these rules, condensed
├── problems/
│   ├── 01-hello-world/
│   │   ├── README.md           # the problem
│   │   ├── LESSON.md           # the teaching
│   │   ├── HINTS.md            # 3 progressive hints, collapsed
│   │   ├── main.go             # starter: compiles, does the wrong thing
│   │   └── main_test.go        # fails until solved
│   ├── 02-.../
│   └── ...
└── capstone/                   # problem 17 lives here (multi-package app)
```

- One directory per problem; each is its own package (mostly `package main`).
- Run one problem: `go test ./problems/01-hello-world/...` and `go run ./problems/01-hello-world`.
- **Starter code must compile.** Functions return zero values or `nil` and carry `// TODO` comments. Tests fail on assertions, never on compile errors.
- Put logic in testable functions; `main()` only calls them and prints. This holds from problem 01 onward (`func Greeting() string`).

## Curriculum (one new concept per problem)

| # | Slug | New concept | Problem idea (domain stays light: loans/listings where natural) | Size |
|---|---|---|---|---|
| 01 | hello-world | `package main`, `func main`, `fmt`, `go mod init`, `go run`, `go test` | Return and print "Hello, World!" | XS |
| 02 | variables-if | `var` vs `:=`, basic types, zero values, `if`/`else`, `switch` (no fall-through) | Classify a credit score into a tier | XS |
| 03 | loops | `for` as the only loop (C-style, while-style, `range`) | FizzBuzz variant: monthly payment schedule labels | XS |
| 04 | functions | Functions in the same file, multiple return values, named returns, variadics | Monthly payment calculator returning payment + total interest | S |
| 05 | errors | `error` as a value, `errors.New`, `fmt.Errorf` + `%w`, `errors.Is/As`; no exceptions | Validate loan inputs, return typed errors | S |
| 06 | slices-maps | Slices (len/cap/append), maps, `range` over both, comma-ok idiom | Group listings by ZIP, compute averages | S |
| 07 | pointers | `&`, `*`, pass-by-value semantics, when to use a pointer | Mutate vs copy a loan record; prove it with tests | S |
| 08 | structs-methods | Structs instead of classes, methods, value vs pointer receivers, `NewX` constructors | `Loan` type with methods (`Balance`, `ApplyPayment`) | M |
| 09 | packages-modules | Multiple files, sub-packages, exported vs unexported (capitalization), `internal/` | Split problem 08 into `loan/` package + `main` | M |
| 10 | interfaces | Implicit interfaces, small interfaces, embedding/composition instead of inheritance | `Notifier` interface with email/SMS fakes | M |
| 11 | table-tests | Table-driven tests, subtests `t.Run`, `t.Helper`, test fakes; `go test -cover` | Write the tests yourself for a provided function | M |
| 12 | json-generics | `encoding/json`, struct tags, generics basics (`[T any]`, constraints) | Parse a JSON loan file; generic `Filter[T]` | M |
| 13 | http-api | `net/http` server, method+path routing (Go 1.22+ `mux`), handlers, JSON responses, `httptest` | REST: `GET/POST /loans` in memory | L |
| 14 | concurrency | Goroutines, channels, `sync.WaitGroup`, `context` cancellation/timeouts | Fan out rate quotes from 3 fake lenders, take fastest, time out | L |
| 15 | postgres | `database/sql` or `pgx`, migrations (plain SQL files), Docker Compose Postgres, Testcontainers-go integration test | Persist loans from problem 13 | L |
| 16 | graphql | `gqlgen` schema-first codegen, resolvers | GraphQL `loan(id)` query + `updateLoan` mutation over problem 15's store | L |
| 17 | capstone | Everything | **Loan task rules API** (below) | XL |

Pacing note in the root README: 01-10 run about 15-30 min each; 11-16 run 45-90 min; the capstone may spill past one day. Never pad a problem to fill time.

## File templates

### `README.md` (problem)
```markdown
# NN - Title
**New concept:** one line
**Builds on:** NN, NN

## Task
2-5 sentences. Concrete inputs and outputs.

## Acceptance criteria
| # | Criterion | Checked by |
|---|---|---|
| 1 | ... | `TestX` |

## Run
go test ./problems/NN-slug/...
go run ./problems/NN-slug

## Stretch (optional)
One extra, no new concept.
```

### `LESSON.md`
Teaches the building blocks needed to solve the problem **without solving it**.

```markdown
# Lesson NN - Concept
## Why Go does it this way          (2-4 bullets: the design reason)
## Syntax                           (minimal code blocks)
## C# ↔ Go
| Go | Closest C# | Difference that will bite you |
|---|---|---|
## Worked example                   (different from the problem's task)
## Gotchas for C# developers        (bullets)
## Check yourself                   (2-3 questions; answers in a collapsed block)
## Go deeper                        (links: go.dev/tour, Effective Go, pkg.go.dev)
```

Rules for lessons:
- Every lesson has the C# ↔ Go table. Name the *closest* C# analogue, and when none exists say so (e.g. implicit interfaces, multiple returns vs `out`/tuples, `error` vs exceptions, no inheritance, `defer` vs `using`/`finally`, goroutines vs `Task`, channels vs `Channel<T>`, struct tags vs attributes, `go.mod` vs `.csproj`/NuGet).
- Code over prose. 60-150 lines per lesson. Scannable.
- The worked example must not be paste-able as the answer.
- Use current idiomatic Go (latest stable toolchain; check go.dev before pinning a version).

### `HINTS.md`
Three progressive hints in `<details>` blocks: (1) the direction, (2) the API or function to look at, (3) near-pseudocode. No full solution.

## Tutor rules (put these in `CLAUDE.md`)

| Rule | Detail |
|---|---|
| No solutions unless asked | Give a solution only when Chris explicitly asks **after** attempting. Point to `HINTS.md` first |
| Review his code first | When he shares a solution: say whether it passes, then idiom feedback (what a Go reviewer would flag), then one "how C# would differ" note |
| Correct his explanations in place | When he explains a concept back, keep his wording, correct each bullet in place, and give a worked instance per claim |
| One concept at a time | Don't introduce concepts from later problems in feedback; mention them as "you'll see this in NN" |
| Tests are the contract | Never weaken a test to make his code pass. If a test is wrong, say so and fix the test explicitly |
| Generate in batches | Create 01-05 first; Chris runs them; then 06-10, 11-16, capstone. Each batch: all files present, `go vet ./...` clean, starter `go test` fails only on assertions |

## Capstone spec: loan task rules API (problem 17)

A small workflow engine: when a loan's data changes, rules decide what tasks to create next (an observe → orient → decide → act loop).

| Area | Requirement |
|---|---|
| Domain | `Loan` (id, borrower, amount, status, documents received, appraisal status, ...), `Task` (id, loanID, type, status, createdAt), `Rule` (`func(old, new Loan) []Task` or an interface) |
| Loop | Observe: a loan update arrives. Orient: diff old vs new. Decide: run all rules. Act: persist new tasks, idempotently (no duplicate open task of the same type) |
| Rules (start with 4) | Missing income doc → `RequestIncomeDoc`; appraisal ordered → `ScheduleAppraisal`; amount > threshold → `UnderwriterReview`; status → closing → `PrepareClosingDocs` |
| REST | `POST /loans`, `PATCH /loans/{id}`, `GET /loans/{id}/tasks`, `PATCH /tasks/{id}` |
| GraphQL | `loan(id){ tasks { ... } }`, `updateLoan(input)`, `completeTask(id)` via gqlgen, sharing the same service layer |
| Storage | Postgres, SQL migrations, transactions around "update loan + create tasks" |
| Tests | Table-driven unit tests per rule; Testcontainers-go integration test for the update → tasks flow; `httptest` for one REST route |
| Layout | `cmd/server`, `internal/loan`, `internal/rules`, `internal/store`, `internal/api/rest`, `internal/api/graphql` |
| Ops | `docker compose up` brings up Postgres and the app; `Makefile` with `test`, `run`, `migrate` |
| README | What it does, the loop diagram, how the rules are tested, design trade-offs, what's next |

The capstone README is public-facing portfolio material: describe it as a learning project in Go, with no employer names, and claim nothing it doesn't do.

## Devcontainer

| Item | Choice |
|---|---|
| Base image | `mcr.microsoft.com/devcontainers/go` (latest stable Go tag) |
| Features | `docker-in-docker` (Compose Postgres + Testcontainers-go) |
| VS Code extensions | `golang.go` (gopls, delve) |
| Tools on create | `go install` gopls, dlv; gqlgen is used via `go run github.com/99designs/gqlgen` at problem 16 |

## Anti-patterns (Claude must avoid)

- Writing solutions into starter files or lessons
- Introducing two new concepts in one problem
- Starter code that doesn't compile
- Lessons longer than about 150 lines, or prose where a table or code block would do
- Translating C# patterns literally (getter/setter classes, exception-driven flow, deep inheritance) as if they were idiomatic Go
