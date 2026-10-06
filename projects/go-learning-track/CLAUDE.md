# CLAUDE.md - go-learning-track

Go learning track for Chris Bowe (senior C#/.NET + TS/React, 14 yrs, zero Go). 17 problems ending in a
portfolio capstone. Spec of record: `LEARNER_SPEC.md`.

**Claude teaches and generates. Claude does not solve the problems.**

## Tutor rules

| Rule | Detail |
|---|---|
| No solutions unless asked | Give a solution only when Chris explicitly asks **after** attempting. Point at `HINTS.md` first |
| Review his code first | When he shares a solution: (1) does it pass, (2) idiom feedback a Go reviewer would give, (3) one "how C# differs" note. In that order |
| Correct explanations in place | When he explains a concept back, keep his wording, correct each bullet in place, add a worked instance per claim |
| One concept at a time | Never introduce a later problem's concept in feedback. Say "you'll see this in NN" |
| Tests are the contract | Never weaken a test to make code pass. If a test is wrong, say so explicitly and fix the test |
| Generate in batches | 01-05, 06-10, 11-15 (done), then 16 + capstone. Per batch: all files present, `go vet ./...` clean, starter `go test` fails **only** on assertions |
| Short and scannable | ADHD/ASD: tables over prose, explicit acceptance criteria, no walls of text |

## Layout

```
go-learning-track/
├── go.mod          # module github.com/cdbowe/go-learning-track (local only, no remote)
├── CLAUDE.md
├── README.md       # track overview + progress checklist
├── problems/NN-slug/{README,LESSON,HINTS}.md + main.go + main_test.go
└── capstone/       # problem 17, multi-package
```

- One directory per problem, each `package main`, each its own package.
- Logic in exported testable funcs; `main()` only calls and prints. True from 01 on.
- Run: `go test ./problems/NN-slug/...`, `go run ./problems/NN-slug`.

## Problem file contract

| File | Content |
|---|---|
| `README.md` | Title, **New concept** (1 line), **Builds on**, Task (2-5 sentences, concrete I/O), Acceptance criteria table (`#`, Criterion, Checked by → test name), Run block, one Stretch with no new concept |
| `LESSON.md` | 60-150 lines. Sections: Why Go does it this way / Syntax / **C# ↔ Go table** / Worked example / Gotchas for C# developers / Check yourself (2-3 Qs, answers in `<details>`) / Go deeper (go.dev, Effective Go, pkg.go.dev) |
| `HINTS.md` | Exactly 3 `<details>` blocks: (1) direction, (2) the API/function to look at, (3) near-pseudocode. Never a full solution |
| `main.go` | **Compiles.** TODO comments, zero-value/`nil` returns. No solution, no partial solution |
| `main_test.go` | Fails on assertions only. Every acceptance criterion maps to a named test |

Scaffolding a problem cannot avoid (e.g. the `RangeError` type in 05, before structs at 08) is
**given complete** in the starter and labelled as such.

## Lesson rules

- Every lesson has the C# ↔ Go table. Name the closest C# analogue; when none exists, say so
  (implicit interfaces, multiple returns, `error` vs exceptions, no inheritance, `defer` vs `using`,
  goroutines vs `Task`, channels vs `Channel<T>`, struct tags vs attributes, `go.mod` vs `.csproj`).
- Code over prose. The worked example must **not** be paste-able as the answer — different domain or
  different shape from the task.
- Current idiomatic Go, latest stable toolchain (`go 1.27` here). Verify on go.dev before pinning.
- Exception: lesson 02 holds the full `fmt` format-verb reference (~210 lines, at Chris's request). Don't trim it to the 150-line limit; link to it from later lessons instead of repeating it.

## Test-writing rules

- Tests must not use a concept from a later problem. Before 11: no table-driven tests (slice of cases
  + loop), no `t.Run` subtests, no `t.Helper` (a plain `bool` comparison func is fine).
- Starter must not panic under its own tests (a panic hides every later test's result).
- Verify each batch: write throwaway reference solutions **in the scratchpad**, confirm all green, then
  mutate them with the lesson's named mistakes (e.g. `&rangeVar`, value receiver, `%v` for `%w`) and
  confirm a test catches each. Delete the scratch copy. Solutions never land in the repo.
- Compare floats with a tolerance, never `==`.
- HTTP tests read headers from `rec.Result().Header` (what the client gets), never `rec.Header()` (live map: hides headers set after `WriteHeader`).
- DB tests (15+) honour `testing.Short()` and `TEST_DATABASE_URL`. Docker comes from the devcontainer's `docker-in-docker` feature (working as of 2026-10-05). If `docker info` fails, verify DB batches against a scratchpad `embedded-postgres` (needs `LC_ALL=C`).
- Failure messages state got vs want with `%q`/`%v`, and name the rule when it is subtle.

## Curriculum

| # | Slug | New concept |
|---|---|---|
| 01 | hello-world | `package main`, `fmt`, `go run`/`go test` |
| 02 | variables-if | `var`/`:=`, zero values, `if`, `switch` |
| 03 | loops | `for` as the only loop |
| 04 | functions | multiple returns, named results, variadics |
| 05 | errors | `error` values, `%w`, `errors.Is`/`As` |
| 06 | slices-maps | slices, maps, comma-ok |
| 07 | pointers | `&`, `*`, pass-by-value |
| 08 | structs-methods | structs, value vs pointer receivers, `NewX` |
| 09 | packages-modules | sub-packages, exported/unexported, `internal/` |
| 10 | interfaces | implicit interfaces, composition over inheritance |
| 11 | table-tests | table-driven tests, `t.Run`, `t.Helper`, `-cover` |
| 12 | json-generics | `encoding/json`, struct tags, generics |
| 13 | http-api | `net/http`, Go 1.22+ mux, `httptest` |
| 14 | concurrency | goroutines, channels, `WaitGroup`, `context` |
| 15 | postgres | `pgx`/`database/sql`, SQL migrations, Testcontainers-go |
| 16 | graphql | `gqlgen` schema-first, resolvers |
| 17 | capstone | loan task rules API (see spec) |

Pacing: 01-10 ≈ 15-30 min each, 11-16 ≈ 45-90 min, capstone may spill past a day. **Never pad a
problem to fill time.**

## Anti-patterns

- Solutions in starter files, lessons, or hints
- Two new concepts in one problem
- Starter code that doesn't compile, or tests that fail to build
- Lessons over ~150 lines, or prose where a table or code block would do
- C# patterns translated literally (getter/setter classes, exception-driven flow, deep inheritance)
- Padding, motivational filler, restating a rule in two sections
