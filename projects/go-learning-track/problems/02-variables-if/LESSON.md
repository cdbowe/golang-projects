# Lesson 02 - Declarations, zero values, if, switch

## Why Go does it this way

- **Zero values, not null.** Every type has a usable zero: `0`, `""`, `false`, `nil` for pointers/maps/slices. A declared variable is never undefined, so `NullReferenceException` has no direct equivalent for value types.
- **`:=` is sugar, not dynamic typing.** Types are inferred at compile time and fixed forever.
- **No implicit numeric conversion.** `int` + `int64` won't compile. Go would rather you say which you meant.
- **`switch` does not fall through.** The C bug class is gone; `fallthrough` exists but is rare.

## Syntax

```go
var score int        // 0
var name string      // ""
var ok bool          // false
var rate float64 = 0.065

score := 712         // inferred int, function scope only
rate2 := 0.065       // inferred float64
const MaxScore = 850 // compile-time constant, untyped

if score > 800 {
	// ...
} else if score > 700 {
	// ...
} else {
	// ...
}

if tier := Tier(score); tier == "Poor" { // scoped to the if/else chain
	// tier visible here and in else
}

switch {                 // no condition == switch true
case score >= 800:
	return "Exceptional"
case score >= 740:
	return "Very Good"
default:
	return "Invalid"
}

switch tier {            // switch on a value
case "Fair", "Good":     // several values, one case
	return "Manual review"
}
```

No parentheses around conditions. Braces are mandatory — no single-line `if (x) return y;`.

## Format verbs (`fmt.Printf`, `Sprintf`, `Errorf`)

Go has **no string interpolation**. Build strings with `+`, or with a format string whose `%` verbs are
filled by the arguments in order. Examples are real output, with `loan := Loan{"L-1", 1200.5}`.

### Verbs

| Verb | Applies to | Example arg | Result |
|---|---|---|---|
| `%v` | any — default format | `loan` / `&loan` | `{L-1 1200.5}` / `&{L-1 1200.5}` |
| `%+v` | any — adds struct field names | `loan` | `{ID:L-1 Balance:1200.5}` |
| `%#v` | any — Go-syntax literal | `loan` | `main.Loan{ID:"L-1", Balance:1200.5}` |
| `%T` | any — the type | `loan` | `main.Loan` |
| `%%` | — literal percent sign | | `%` |
| `%t` | bool | `true` | `true` |
| `%d` | integer, base 10 | `42` | `42` |
| `%b` / `%o` / `%O` | integer, base 2 / 8 / 8 with `0o` | `5` / `8` / `8` | `101` / `10` / `0o10` |
| `%x` / `%X` | integer, hex | `255` | `ff` / `FF` |
| `%c` | integer as a character | `65` | `A` |
| `%q` | integer as a quoted character | `65` | `'A'` |
| `%U` | integer as a Unicode code point | `0x1F600` | `U+1F600` |
| `%f` / `%F` | float, decimal | `3.14159` | `3.141590` (6 places by default) |
| `%e` / `%E` | float, scientific | `1234.5678` | `1.234568e+03` / `1.234568E+03` |
| `%g` / `%G` | float, shortest of `%f`/`%e` | `1234.5678` / `1e21` | `1234.5678` / `1e+21` |
| `%x` / `%X` / `%b` | float, hex / binary exponent | `1.0` | `0x1p+00` / `4503599627370496p-52` (rare) |
| `%s` | string, `[]byte`, `error`, anything with `String()` | `"hi"` | `hi` |
| `%q` | string, double-quoted and escaped | `"hi\n"` | `"hi\n"` |
| `%x` / `%X` | string or `[]byte` as hex bytes | `"hi"` | `6869` |
| `%p` | pointer, slice, map, chan, func — address | `&loan` | `0xc000010030` (varies per run) |
| `%w` | `error` — **`fmt.Errorf` only**, wraps the cause (05) | `err` | the error's text |

Slices, arrays and maps apply the verb per element: `%v` of `[]int{1,2}` → `[1 2]`; `%q` of
`[]string{"a","b"}` → `["a" "b"]`; a map prints `map[a:1 b:2]` with keys **sorted**. `nil` prints
`<nil>`; an `error` prints its `Error()` text.

### Flags, width, precision

Form: `%[flags][width][.precision]verb`.

| Spec | Meaning | Example | Result |
|---|---|---|---|
| `%8.2f` | width 8, 2 decimals, right-aligned | `3.14159` | `"    3.14"` |
| `%-8.2f` | `-` left-aligns | `3.14159` | `"3.14    "` |
| `%08.2f` / `%06d` | `0` pads with zeros | `3.14159` / `42` | `00003.14` / `000042` |
| `%+d` | `+` always prints the sign | `42` | `+42` |
| `% d` / `% x` | space: room for a sign / spaces between hex bytes | `42` / `"hi"` | `" 42"` / `68 69` |
| `%#x` / `%#q` / `%#U` | `#` alternate form | `255` / `"hi"` / `65` | `0xff` / `` `hi` `` / `U+0041 'A'` |
| `%10s` / `%-10s` | pad a string to width 10 | `"hi"` | `"        hi"` / `"hi        "` |
| `%.3s` | precision on a string truncates | `"abcdef"` | `abc` |
| `%*d` | width taken from an argument | `5, 42` | `"   42"` |
| `%[2]d %[1]d` | explicit argument index | `1, 2` | `2 1` |
| `%[1]d %[1]q` | reuse one argument | `65` | `65 'A'` |

### Which function

| Function | Goes to | Adds spaces / newline |
|---|---|---|
| `Print` / `Println` | stdout | `Println`: spaces between args + `\n`. **Ignores `%` verbs** |
| `Printf` | stdout | Neither — end the format with `\n` yourself |
| `Sprint` / `Sprintln` / `Sprintf` | returns a `string` | same rules as above |
| `Fprint` / `Fprintln` / `Fprintf` | any writer, e.g. `&strings.Builder` (03) | same rules as above |
| `Errorf` | returns an `error` | Supports `%w` (several allowed since Go 1.20) |

### Mistakes are printed, not thrown

| Mistake | Output | `go vet` says |
|---|---|---|
| wrong type | `Sprintf("%d", "hi")` → `%!d(string=hi)` | `format %d has arg "hi" of wrong type string` |
| missing arg | `Sprintf("%d %d", 1)` → `1 %!d(MISSING)` | `format %d reads arg #2, but call has 1 arg` |
| extra arg | `Sprintf("%d", 1, 2)` → `1%!(EXTRA int=2)` | `call needs 1 arg but has 2 args` |
| `%w` outside `Errorf` | `%!w(*errors.errorString=&{x})` | `fmt.Sprintf does not support error-wrapping directive %w` |
| verb in `Println` | `Println("%.2f", x)` prints the `%` literally | `Println call has possible Printf formatting directive %.2f` |

`go vet` only checks constant format strings — another reason to keep them literal.

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `var x int` | `int x;` | Go's is already initialized to `0`; C#'s is "unassigned" and the compiler stops you from reading it |
| `x := 5` | `var x = 5` | Go's `:=` only works inside a function and only when at least one name on the left is new |
| zero value | `default(T)` | You get it automatically on declaration, not just via `default` |
| `""` | `null` / `string.Empty` | A Go `string` is never null. There is no `string?` |
| `switch` | `switch` | No `break` needed, no implicit fall-through, and `case` can be any expression, not just a constant |
| `switch { case cond: }` | `if`/`else if` chain | Idiomatic Go prefers the conditionless switch where C# would chain `else if` |
| `const MaxScore = 850` | `const int MaxScore = 850` | A Go constant can be *untyped* and adapt to the context it's used in (`int`, `float64`, ...) |
| `int` | `int` (32-bit) | Go's `int` is 64-bit on this machine. Convert explicitly: `float64(score)`, `int(rate)` |
| `fmt.Sprintf("%-10s %8.2f", n, x)` | `$"{n,-10} {x,8:F2}"` | No interpolation: arguments are positional. A type mismatch is not a compile error — it prints `%!d(...)`; `go vet` catches it |

## Worked example

Classify an LTV ratio — the same shape as the task, different domain:

```go
// LTVBand labels a loan-to-value ratio. Input is a fraction: 0.8 == 80%.
func LTVBand(ltv float64) string {
	if ltv <= 0 || ltv > 1.25 {
		return "Invalid" // guard clauses first: cheap exits before real logic
	}

	switch {
	case ltv <= 0.60:
		return "Low"
	case ltv <= 0.80:
		return "Standard"
	case ltv <= 0.95:
		return "Elevated"
	default:
		return "High"
	}
}

// Requires reports whether an LTV band needs extra paperwork.
func Requires(band string) string {
	switch band {
	case "Low", "Standard":
		return "none"
	case "Elevated", "High":
		return "PMI"
	default:
		return "Unknown"
	}
}
```

Two patterns worth copying: **guard clause first**, then ordered `case`s from one end of the range to
the other so each case only needs one comparison.

## Gotchas for C# developers

- `:=` at package level is a compile error. Outside functions, use `var` or `const`.
- Re-declaring with `:=` in the same scope fails (`no new variables on left side of :=`); use `=`.
- An unused local variable is a *compile error*. An unused package-level `var` is fine.
- `score := 712` then `score = "high"` fails: the type was fixed at declaration.
- Shadowing is legal and silent: `if tier := ...` inside a function that already has `tier` creates a
  second variable. `go vet` won't flag it; read your braces.
- `string` comparison with `==` compares contents, like C#. No `.Equals` needed.
- Switching on a value requires the same type everywhere: `switch score { case "700": }` won't compile.

## Check yourself

1. `var count int` then `fmt.Println(count)` — what prints, and what would C# do with the equivalent?
2. Why does `x := 5; x := 6` fail but `x := 5; x, y := 6, 7` compile?
3. When would you pick `switch {}` over an `if`/`else if` chain?

<details>
<summary>Answers</summary>

1. Go prints `0`. C# refuses to compile: "use of unassigned local variable". Go's declaration *is* an
   initialization to the zero value.
2. `:=` requires at least one new name on the left. In the second case `y` is new, so it's allowed —
   and `x` is just assigned, not redeclared.
3. When each branch tests the same subject and you want them read as a flat list of alternatives. It
   also makes `default` explicit, which the compiler won't force but reviewers will ask for.

</details>

## Go deeper

- [pkg.go.dev/fmt](https://pkg.go.dev/fmt) — the authoritative verb list
- [Tour of Go: flow control](https://go.dev/tour/flowcontrol/1)
- [Tour of Go: zero values](https://go.dev/tour/basics/12)
- [Effective Go: control structures](https://go.dev/doc/effective_go#control-structures)
