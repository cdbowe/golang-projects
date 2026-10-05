# Lesson 11 - Table-driven tests

## Why Go does it this way

- **A test case is data.** Inputs and expected outputs go in a slice of structs; one loop runs them. Adding a case is adding a line, not a method.
- **Subtests give each row a name.** `t.Run(name, func(t *testing.T) {...})` reports, filters, and fails per case — `TestQuote/credit_619` is addressable from the command line.
- **No assertion or mocking library needed.** `if got != want { t.Errorf(...) }` plus a 5-line fake struct covers almost everything. `t.Helper()` keeps failure lines pointing at the caller.
- **Coverage is built in.** `go test -cover` — no Coverlet, no collector config.

## Syntax

```go
func TestTier(t *testing.T) {
	tests := []struct {
		name  string
		score int
		want  string
	}{
		{"floor of Poor", 300, "Poor"},
		{"just below Fair", 579, "Poor"},
		{"zero value", 0, "Invalid"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { // inner t: failures belong to this subtest
			if got := Tier(tc.score); got != tc.want {
				t.Errorf("Tier(%d) = %q, want %q", tc.score, got, tc.want)
			}
		})
	}
}
```

A helper:

```go
func assertNear(t *testing.T, got, want float64) {
	t.Helper() // failure is reported at the caller's line, not here
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("got %v, want %v", got, want)
	}
}
```

A fake with configurable behaviour, and a function value inside the table:

```go
type stubSource struct {
	rate float64
	err  error
}

func (s stubSource) BaseRate(int) (float64, error) { return s.rate, s.err }

tests := []struct {
	name    string
	src     RateSource  // each row can carry its own fake
	wantErr error       // nil means "expect success"
}{...}
```

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| slice-of-struct table | `[Theory]` + `[InlineData]` / `[MemberData]` | Plain Go code, no attributes. Any type can go in a row, including fakes and funcs |
| `t.Run(name, fn)` | one theory row in Test Explorer | You choose the name. Spaces become `_` in output and in `-run` patterns |
| `-run 'TestQuote/credit'` | `--filter` | Regex per path segment, split on `/` |
| `t.Errorf` / `t.Fatalf` | `Assert.Equal` / throw | `Errorf` keeps going; `Fatalf` stops **this subtest only** — the loop continues to the next row |
| `t.Helper()` | `[StackTraceHidden]` | Marks the function so failure lines point to the caller |
| hand-written fake | Moq / NSubstitute | A struct with a method satisfies the interface implicitly (10). No setup DSL |
| `go test -cover` | Coverlet | Built in. `-coverprofile` + `go tool cover -func` / `-html` for detail |
| `t.Log` | `ITestOutputHelper` | Shown only for failures, or always with `-v` |
| mutation testing (`mutants.sh`) | Stryker.NET | Same idea, hand-rolled here with build tags |

## Worked example

Problem 04's `TotalFees`, as a table — note the variadic slice in each row:

```go
func TestTotalFees(t *testing.T) {
	tests := []struct {
		name  string
		base  float64
		extra []float64
		want  float64
	}{
		{"no extras", 1200, nil, 1200},
		{"two extras", 1200, []float64{450, 85.50}, 1735.50},
		{"negative extra", 100, []float64{-25}, 75},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := TotalFees(tc.base, tc.extra...)
			assertNear(t, got, tc.want)
		})
	}
}
```

For error cases, give the row a `wantErr error` and branch inside the subtest:

```go
if tc.wantErr != nil {
	if !errors.Is(err, tc.wantErr) {
		t.Fatalf("err = %v, want %v", err, tc.wantErr)
	}
	return // nothing else to check for this row
}
if err != nil {
	t.Fatalf("unexpected error: %v", err)
}
```

## Gotchas for C# developers

- Name every row. `t.Run("", ...)` and `t.Run(fmt.Sprint(i), ...)` make failures unreadable.
- Use the subtest's `t`, not the outer one. Calling the outer `t.Fatal` inside `t.Run` misreports the failure.
- Test **both sides of every boundary**: 619 and 620, 0.80 and 0.81. Mutation testing exists because "a test for each branch" isn't the same as "a test for each boundary".
- 100% coverage doesn't mean correct. Every line can run while an off-by-one survives — that's what the mutants show.
- A fake that ignores its input can't catch a bug that passes the wrong input. If `Quote` must ask for the *right* term, your fake must answer differently per term.
- Floats need a tolerance: `0.0625 + 0.0025` is not exactly `0.065`.

## Check yourself

1. What's the difference between `t.Fatalf` inside a subtest and inside a plain loop without `t.Run`?
2. Your suite has 100% coverage but `mutant_ltv_surcharge` survives. What's missing?
3. How do you check that `Quote` still wraps the source's own error?

<details>
<summary>Answers</summary>

1. In a subtest it stops that row only; the loop runs the remaining rows. In a plain loop it stops the
   whole test at the first failing row.
2. A case at exactly `ltv = 0.80`. Coverage only needs *some* LTV above and below; the boundary needs
   the exact value, where `>` and `>=` disagree.
3. Make the fake return a sentinel you own (`errSourceDown := errors.New("source down")`), then assert
   `errors.Is(err, errSourceDown)` as well as `errors.Is(err, ErrRateUnavailable)`.

</details>

## Go deeper

- [Go wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests)
- [Go blog: Using Subtests and Sub-benchmarks](https://go.dev/blog/subtests)
- [pkg.go.dev/testing](https://pkg.go.dev/testing) — `T.Run`, `T.Helper`, `T.Cleanup`
- [Go blog: The cover story](https://go.dev/blog/cover)
