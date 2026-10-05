# Lesson 01 - Packages, main, and the test loop

## Why Go does it this way

- **One binary, no runtime.** `go build` emits a single static executable. No CLR, no `dotnet` on the target box.
- **The compiler is the linter.** Unused imports and unused local variables are *compile errors*, not warnings. Go would rather fail now than let rot accumulate.
- **Testing ships in the standard library.** No xUnit, no NUnit, no package reference. `testing` + `go test` is the whole story.
- **Convention over configuration.** File named `*_test.go`, function named `TestXxx(t *testing.T)`, directory == package. There is no project file listing your sources.

## Syntax

```go
package main // an executable. Any other name == a library.

import "fmt" // one import

import (     // several
	"fmt"
	"strings"
)

func Greeting() string { // Exported: capital G. Returns string.
	return "hi"
}

func main() { // the entry point: no args, no return
	fmt.Println(Greeting())
}
```

Test file, `main_test.go`, same package:

```go
package main

import "testing"

func TestThing(t *testing.T) {
	if got := Greeting(); got != "hi" {
		t.Errorf("Greeting() = %q, want %q", got, "hi")
	}
}
```

There is no `Assert.Equal`. You write the `if`, then call `t.Errorf` (fail, keep going) or `t.Fatalf`
(fail, stop this test).

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `package main` + `func main()` | `class Program { static void Main() }` | No class. `main` takes no args — CLI args come from `os.Args` |
| `go.mod` | `.csproj` + NuGet | One `go.mod` per module, not per project. No `<Compile Include>`: every `.go` file in the folder is in the build |
| directory = package | namespace | A namespace can span folders; a Go package **is** one folder. Nesting folders creates unrelated packages, not parent/child |
| `Greeting` vs `greeting` | `public` vs `internal` | Capitalization *is* the access modifier. No keyword exists |
| `go test ./...` | `dotnet test` | Tests live beside the code in the same package, not in a separate test project |
| `fmt.Println` | `Console.WriteLine` | `fmt` verbs (`%q`, `%v`, `%d`) replace interpolation in test messages |
| `t.Errorf` | `Assert.Equal` | No assertion library in stdlib. The `if` is the assertion |

## Worked example

A different function, wired the same way:

```go
package main

import "fmt"

// Initials turns "Chris Bowe" into "CB". Logic lives here, so it is testable.
func Initials(first, last string) string {
	return string(first[0]) + string(last[0])
}

func main() {
	fmt.Println(Initials("Chris", "Bowe")) // CB
}
```

```go
func TestInitials(t *testing.T) {
	if got := Initials("Chris", "Bowe"); got != "CB" {
		t.Errorf("Initials() = %q, want %q", got, "CB")
	}
}
```

Note `first, last string` — the shared type is written once. And `%q` quotes the value in the failure
message, which saves you when the bug is a trailing space.

## Gotchas for C# developers

- `func main()` returning a value won't compile. Exit codes come from `os.Exit(1)`.
- `import "strings"` without using `strings` won't compile. Delete it or use it.
- `gofmt` is not negotiable — tabs, brace placement, import order. Run `gofmt -l .`; your editor with
  `golang.go` does it on save. Nobody argues about style in Go reviews.
- `go test ./problems/01-hello-world` with no `/...` works here, but `/...` means "this and everything
  below", which is what you want once there are sub-packages.
- Verbose output: `go test -v ./...`. Single test: `go test -run TestGreetingText ./...`.

## Check yourself

1. Why can a Go file not be in two packages at once?
2. What happens if you leave `import "strings"` in a file that never calls it?
3. Where does the test for `problems/01-hello-world/main.go` live, and in what package?

<details>
<summary>Answers</summary>

1. Package identity comes from the directory. Every `.go` file in a folder declares the same package
   name, so the file's location decides its package. (C#'s `namespace` is a declaration you can put
   anywhere; Go's is the filesystem.)
2. Compile error: `"strings" imported and not used`. Same for an unused *local* variable. Unused
   struct fields and unused exported functions are fine.
3. `problems/01-hello-world/main_test.go`, `package main` — same folder, same package, so it can call
   the unexported identifiers too.

</details>

## Go deeper

- [A Tour of Go](https://go.dev/tour/welcome/1) — packages, functions (first 10 slides)
- [Effective Go: names](https://go.dev/doc/effective_go#names)
- [pkg.go.dev/fmt](https://pkg.go.dev/fmt) — the verb table is worth a bookmark
- [pkg.go.dev/testing](https://pkg.go.dev/testing)
