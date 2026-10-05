# 11 - Test the rate quote

**New concept:** table-driven tests — `t.Run` subtests, `t.Helper`, hand-written fakes, `go test -cover`
**Builds on:** 05, 06, 10

## Task

This time the code is written and **you write the tests.** `Quote` in `quote.go` is complete and
correct; its doc comment is the spec. Replace the placeholder `TestQuote` in `main_test.go` with a
table-driven test that pins down every rule.

Your tests must:

- be **one table** of cases, each run with `t.Run(tc.name, ...)`
- use **hand-written fakes** for `RateSource` (no mocking library): at least one that returns rates per
  term, and one that returns an error
- use a **`t.Helper()`** assertion function for comparing rates (floats need a tolerance)
- check errors with **`errors.Is`**, including the *source's own* error through the wrap

`mutant_*.go` are deliberately buggy copies of `Quote`. `mutants.sh` runs your tests against each one;
a good suite fails on every mutant. Don't open them first — they give away which cases matter.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | `TestQuote` passes against the real `Quote` | `go test ./problems/11-table-tests/...` |
| 2 | One case table, each case a named `t.Run` subtest | review; `go test -v` lists `TestQuote/<name>` |
| 3 | A `t.Helper()` function does the float comparison | review; failures point at the case, not the helper |
| 4 | Fakes are plain types implementing `RateSource` | review |
| 5 | 100% statement coverage of `Quote` | `go tool cover -func` shows `Quote 100.0%` |
| 6 | All 8 mutants caught | `bash problems/11-table-tests/mutants.sh` prints `8/8` |

## Run

```bash
go test -v ./problems/11-table-tests/...
go test -run 'TestQuote/credit' -v ./problems/11-table-tests/...      # only matching subtests
go test -coverprofile=/tmp/c.out ./problems/11-table-tests/... && go tool cover -func=/tmp/c.out
bash problems/11-table-tests/mutants.sh
```

Overall package coverage will be below 100% because `main()` and the demo `todaysRates` aren't tested.
Criterion 5 is about the `Quote` line only.

## Stretch (optional)

Rewrite problem 02's `Tier` tests as a table in under 30 lines. Same assertions, a fraction of the
code.
