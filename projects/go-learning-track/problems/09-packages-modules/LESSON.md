# Lesson 09 - Packages, files, and visibility

## Why Go does it this way

- **A package is a directory.** Every `.go` file in a folder is one package; they share all names, unexported included. Splitting a file costs nothing and changes nothing semantically.
- **Capitalization is the only access modifier.** `New`, `Loan`, `Balance` are exported; `validID`, `principal` are not. It's visible at every call site — `loan.New` is obviously public, `validID` obviously isn't.
- **The package name is part of every call**, so names are chosen to read well *qualified*: `loan.New`, `http.Get`, `strings.TrimSpace`. Repeating the package name (`loan.NewLoan`) is called stutter, and reviewers flag it.
- **`internal/` is enforced by the compiler.** Code under `a/b/internal/x` can be imported only from within `a/b/...`. It's how you share code inside a module without making it public API.

## Syntax

```go
// internal/rates/rates.go
package rates                    // name = last path element, by convention

const maxRate = 0.25             // unexported: package-wide, invisible outside

// Monthly converts an annual rate. Exported: comment starts with its name.
func Monthly(annual float64) float64 { return clamp(annual) / 12 }

func clamp(r float64) float64 { return min(max(r, 0), maxRate) } // unexported helper
```

```go
// main.go
package main

import (
	"fmt"                                                      // stdlib first

	"github.com/cdbowe/go-learning-track/some/path/internal/rates" // module path + directory
)

func main() { fmt.Println(rates.Monthly(0.065)) }   // rates.clamp → compile error
```

Import path = **module path from `go.mod`** + **directory**. The identifier you use (`rates`) comes from
the `package` clause, not the folder name.

Two test styles for one package:

```go
package rates      // rates_internal_test.go: white-box, can call clamp()
package rates_test // rates_test.go: black-box, must import ".../rates" and use rates.Monthly
```

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| package (directory) | namespace + assembly | One folder = one package. No nesting relationship: `loan` and `loan/fees` are unrelated packages |
| `Exported` / `unexported` | `public` / `internal` | Unexported is *package*-scoped, not file- or type-scoped. No `private`, `protected`, `friend` |
| `internal/` directory | `internal` + `InternalsVisibleTo` | Scoped by directory tree, enforced at compile time, no attribute needed |
| module (`go.mod`) | solution / `.csproj` | One module, many packages; no project file lists sources |
| import path | `using` + project reference | The import *is* the reference. No separate "add reference" step for code in your module |
| `package x_test` | separate test project | Lives in the same folder, compiled as a separate package |
| `go doc ./path` | XML docs / IntelliSense | Reads plain `//` comments above exported names. No `<summary>` tags |
| no circular imports | circular project refs (also banned) | Go rejects package-level cycles, not just assembly-level ones. Design dependencies one way |

## Worked example

A `fees` package split across two files, imported by `main`:

```
fees/
├── fees.go      // package fees: Total(base float64, extra ...float64) float64
└── round.go     // package fees: func roundCents(v float64) float64  (unexported)
```

```go
// fees/round.go
package fees

import "math"

func roundCents(v float64) float64 { return math.Round(v*100) / 100 }
```

```go
// fees/fees.go
package fees

// Total sums the fees and rounds to the cent.
func Total(base float64, extra ...float64) float64 {
	t := base
	for _, e := range extra {
		t += e
	}
	return roundCents(t) // same package, different file: no import needed
}
```

```go
// main.go — fees.Total works, fees.roundCents is a compile error:
// "name roundCents not exported by package fees"
fmt.Println(fees.Total(1200, 450, 85.505))
```

## Gotchas for C# developers

- Name collisions: a local variable named `loan` shadows the imported package `loan`. Call it `l`.
- Every file in a folder needs the **same** `package` line (except `_test` packages). A mismatch is a build error.
- Unused imports are compile errors here too — moving code between packages tends to leave one behind.
- An exported func returning an unexported type compiles, but callers can't name it. Reviewers flag it.
- `go doc` shows only exported names. If something isn't there, it isn't your API.
- Doc comments start with the name: `// New validates...`. Linters check this.
- Don't create `utils`/`common`/`helpers` packages. Name packages for what they provide.

## Check yourself

1. Why is `loan.New` preferred over `loan.NewLoan`?
2. Can `loan_test.go` (package `loan_test`) call `validID`? Can `validate_test.go`?
3. What stops `problems/10-interfaces` from importing this `internal/loan`?

<details>
<summary>Answers</summary>

1. Callers always write the package name, so `loan.NewLoan` says "loan" twice. When a package has one
   main type, `New` is the convention (`list.New`, `errors.New`).
2. `loan_test.go`: no — it's a separate package and `validID` is unexported. `validate_test.go`: yes —
   it declares `package loan`.
3. The compiler: `internal/` under `problems/09-packages-modules` is importable only from inside
   `problems/09-packages-modules/...`. The error is `use of internal package ... not allowed`.

</details>

## Go deeper

- [How to Write Go Code: packages and modules](https://go.dev/doc/code)
- [Go blog: Package names](https://go.dev/blog/package-names) — short, and the stutter rule's source
- [Go Code Review Comments: doc comments](https://go.dev/wiki/CodeReviewComments#doc-comments)
