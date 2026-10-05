# Lesson 05 - Errors are values

## Why Go does it this way

- **No exceptions, no stack unwinding.** A function that can fail returns `error` as its last result. The failure path is in the signature, so you can see it at the call site without reading the body.
- **`error` is just an interface** with one method: `Error() string`. Anything implementing it is an error — no base class to inherit.
- **Wrapping is explicit.** `%w` records a cause chain. `errors.Is` walks the chain for identity; `errors.As` walks it looking for a type.
- **`panic` exists but isn't control flow.** It's for "this program is broken" (nil map write, index out of range), not "the user typed a bad amount."

## Syntax

```go
// Sentinel: one value, created once, compared by identity.
var ErrNotFound = errors.New("not found")

// Returning
func Find(id string) (string, error) {
	if id == "" {
		return "", ErrNotFound
	}
	return "loan-" + id, nil // nil error == success
}

// Checking — the shape you'll type a thousand times
name, err := Find("")
if err != nil {
	return err
}

// Adding context while keeping the cause
if err != nil {
	return fmt.Errorf("loading loan %s: %w", id, err) // %w wraps; %v flattens
}

// Identity: is this (anywhere in the chain) that sentinel?
if errors.Is(err, ErrNotFound) { }

// Type: is there a *RangeError in the chain, and what's in it?
var re *RangeError
if errors.As(err, &re) {
	fmt.Println(re.Field, re.Min, re.Max)
}
```

A custom error type is any type with an `Error() string` method:

```go
type RangeError struct{ Field string }

func (e *RangeError) Error() string { return e.Field + " out of range" }

return &RangeError{Field: "amount"} // note the &: the method is on the pointer
```

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `return nil, err` | `throw new Exception()` | No unwinding. Forget `if err != nil` and execution *continues* with a zero value |
| `error` interface | `Exception` base class | No type hierarchy, no stack trace, no message/inner-exception properties. Just `Error() string` |
| `fmt.Errorf("...: %w", err)` | `throw new X("...", inner)` | You choose per call whether to keep the cause (`%w`) or drop it (`%v`) |
| `errors.Is(err, ErrX)` | `catch (XException)` on a singleton | Compares *identity* against a sentinel value, not a type |
| `errors.As(err, &target)` | `catch (RangeException re)` | Needs a pointer to your target variable; returns `bool` and fills it in |
| `panic` / `recover` | `throw` / `catch` | Reserved for bugs and unrecoverable states. Never for validation |
| `defer` | `finally` / `using` | Runs on the way out, including during a panic. You'll use it from 13 on |
| no checked exceptions | — | The compiler won't force you to handle `error`, but `go vet` flags a *dropped* one |
| error strings lowercase, no punctuation | `"Amount is invalid."` | Convention: `"amount out of range"`, because errors get wrapped into longer sentences |

## Worked example

Sentinel + typed error + wrapping, in a different domain:

```go
var ErrNoDocument = errors.New("document missing")

type StaleError struct {
	Name    string
	DaysOld int
}

func (e *StaleError) Error() string {
	return fmt.Sprintf("%s is %d days old", e.Name, e.DaysOld)
}

func CheckDocument(name string, daysOld int) error {
	if name == "" {
		return ErrNoDocument                              // bare sentinel
	}
	if daysOld > 90 {
		return &StaleError{Name: name, DaysOld: daysOld}  // typed, with data
	}
	return nil
}

func CheckFile(name string, daysOld int) error {
	if err := CheckDocument(name, daysOld); err != nil {
		return fmt.Errorf("document check failed: %w", err) // context + cause
	}
	return nil
}
```

What the caller can do with it:

```go
err := CheckFile("paystub.pdf", 120)
// err.Error() == "document check failed: paystub.pdf is 120 days old"

var se *StaleError
if errors.As(err, &se) && se.DaysOld > 365 {
	// errors.As found the wrapped type and gave us its fields
}
errors.Is(err, ErrNoDocument) // false here — the chain holds a *StaleError
```

Note the asymmetry: `errors.Is` for "which error is it", `errors.As` for "give me its data".

## Gotchas for C# developers

- The nil check is `if err != nil`, and it comes *immediately* after the call. Early return, no `else`.
- `%w` vs `%v`: both print the same text. Only `%w` keeps `errors.Is`/`errors.As` working. Nothing warns you.
- `errors.As` needs a **pointer to the target**: `var re *RangeError; errors.As(err, &re)`. Passing `re` panics.
- Methods declared on `*RangeError` mean only `&RangeError{...}` satisfies `error` — a bare `RangeError{}` won't compile as a return value.
- Don't `errors.New` inside a function if callers need to compare it: each call makes a distinct value and `errors.Is` will never match. Sentinels live at package level.
- A non-nil error with a zero-value result is normal. Return `"", err`, not `"", nil`.
- Comparing `err.Error() == "some text"` is a trap — wrapping changes the text. Compare with `errors.Is`.

## Check yourself

1. You wrap with `%v` instead of `%w`. Which of `errors.Is` and `errors.As` still work?
2. Why is `var ErrX = errors.New("x")` at package level, rather than created where it's returned?
3. `errors.As(err, re)` instead of `errors.As(err, &re)` — compile error or runtime panic?

<details>
<summary>Answers</summary>

1. Neither. `%v` flattens the cause into text, so the chain ends there. The message looks identical,
   which is why this bug survives code review.
2. `errors.New` returns a pointer to a new value every call, and `errors.Is` compares identity. Two
   separately created errors with the same text are different errors.
3. Runtime panic: `errors.As: target must be a non-nil pointer`. The parameter is `any`, so the
   compiler can't catch it.

</details>

## Go deeper

- [Go blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors) - the one to read
- [Effective Go: errors](https://go.dev/doc/effective_go#errors)
- [pkg.go.dev/errors](https://pkg.go.dev/errors)
