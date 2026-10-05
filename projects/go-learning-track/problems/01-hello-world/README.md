# 01 - Hello, World!

**New concept:** `package main` + `func main` + `fmt`, and the `go run` / `go test` loop
**Builds on:** nothing

## Task

`Greeting()` returns the string `Hello, World!`. `main()` prints it with `fmt.Println`.

The starter compiles and returns `""`, so the test fails on an assertion, not on a build error. That
split — logic in a testable function, `main()` only wiring — holds for every problem in this track.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | `Greeting()` returns exactly `Hello, World!` | `TestGreetingText` |
| 2 | `Greeting()` returns a non-empty string | `TestGreetingNotEmpty` |
| 3 | `go run ./problems/01-hello-world` prints the greeting | manual |

## Run

```bash
go test ./problems/01-hello-world/...
go run ./problems/01-hello-world
```

## Stretch (optional)

Add `GreetingFor(name string) string` returning `Hello, Chris!` for `"Chris"`. No test provided —
write your own assertion in `main_test.go`.
