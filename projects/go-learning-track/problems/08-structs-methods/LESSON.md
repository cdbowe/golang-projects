# Lesson 08 - Structs and methods

## Why Go does it this way

- **No classes.** A struct holds data; methods are functions with a *receiver*, declared outside the type body. Data and behaviour are joined by a name, not by braces.
- **No constructors.** A `NewX` function is a plain function convention. Make the zero value useful where you can; use `NewX` when there are invariants to check.
- **The receiver is a parameter**, so the pass-by-value rules from 07 apply: a value receiver gets a copy, a pointer receiver can mutate.
- **Encapsulation is per package**, not per type. Lowercase fields are hidden from *other packages* (you'll enforce that in 09); code in the same package can still see them.

## Syntax

```go
// Struct names: 
// 	capitalized == package-public, exported
// 	uncapitalized == package-private, unexported
// Same applies to type, field, and function names

type Account struct {
	owner   string    // unexported: package-private
	limit   float64
	charges []float64 // zero value nil: ready to append
}

func NewAccount(owner string, limit float64) (*Account, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be positive")
	}
	return &Account{owner: owner, limit: limit}, nil // named fields; unlisted ones get zero values
}

func (a *Account) Charge(amt float64) { a.charges = append(a.charges, amt) } // pointer receiver: mutates

func (a *Account) Owner() string { return a.owner } // getter: Owner(), not GetOwner()

a, _ := NewAccount("Chris", 500)
a.Charge(20) // (*a).Charge — Go takes &/deref for you on addressable values
```

Value receiver, for small immutable types:

```go
type Rate float64

func (r Rate) Monthly() Rate { return r / 12 } // methods work on any named type, not just structs
```

## C# ↔ Go

| Go | Closest C# | Difference that will bite you |
|---|---|---|
| `type Loan struct{...}` | `class Loan` | No inheritance, no `virtual`, no access modifiers per member — just capitalization |
| `func (l *Loan) Pay()` | instance method | Declared outside the type. The receiver name (`l`) replaces `this`; keep it short, never `this`/`self` |
| value receiver `(l Loan)` | method on a `struct` | Operates on a copy. Mutations vanish silently |
| `NewLoan(...)` | constructor | Just a function. Nothing stops `Loan{}` being created directly — design the zero value accordingly |
| `ID()` | `public string Id { get; }` | No properties. Getters drop the `Get` prefix; setters are rare — prefer methods that do real work |
| unexported fields | `private` | Visible to the whole *package*, not just the type |
| `Loan{}` zero value | `new Loan()` with defaults | Every field zeroed; no field initializers, no default constructor logic |
| struct `==` | record equality | Works only if every field is comparable — a slice field makes `==` a compile error |
| methods on `type Rate float64` | extension methods | Real methods, but only on types declared in your own package |

## Worked example

```go
type Escrow struct {
	required float64
	held     []float64
}

func NewEscrow(required float64) *Escrow { // can't fail: no error result
	return &Escrow{required: required}
}

func (e *Escrow) Deposit(amt float64) { e.held = append(e.held, amt) }

func (e *Escrow) Shortfall() float64 {
	total := 0.0
	for _, h := range e.held {
		total += h
	}
	return max(e.required-total, 0)
}

func (e *Escrow) Deposits() []float64 { return slices.Clone(e.held) } // don't leak internal state
```

Every method uses `*Escrow`, even `Shortfall`, which only reads. That's the rule below.

## Choosing a receiver

| Use pointer `*T` when | Use value `T` when |
|---|---|
| any method mutates | the type is small and never mutated (`time.Time`, `Rate`) |
| the struct is large | it's a map, func, or chan type |
| it holds a slice/map you'll append/assign | you want copies to be independent by design |

**Be consistent:** if one method needs `*T`, give them all `*T`. Mixed receivers confuse readers and
complicate interface satisfaction (that's 10).

## Gotchas for C# developers

- A value receiver on a mutating method compiles and silently does nothing. `TestApplyPaymentPersists` exists for this.
- Returning an internal slice (`return l.payments`) hands out write access. Clone it.
- `NewLoan` returning `(nil, err)` vs `(&Loan{}, err)`: return `nil` on error so callers can't use a half-built loan.
- Struct literals without field names (`Loan{"L1", 1000, nil}`) break when fields are added. Use names.
- There's no `this.` — and no implicit field access. Inside a method, it's always `l.balance`.
- A method on a nil `*Loan` is callable. It only panics when it touches a field.

## Check yourself

1. `func (l Loan) Pay(x float64) { l.paid += x }` — compile error, runtime error, or silent bug?
2. Why is `Payments()` returning `slices.Clone(...)` and not the field itself?
3. Store the total paid, or recompute it from the payment history each time?

<details>
<summary>Answers</summary>

1. Silent bug. It compiles and runs; `l` is a copy, so the caller's loan never changes.
2. A slice header shares its backing array (06). Returning the field lets callers rewrite history and
   desync `Balance()`.
3. Either passes the tests. Recomputing has one source of truth; a running total is O(1) but must be
   kept in sync. Pick one and be ready to say why — reviewers ask.

</details>

## Go deeper

- [Tour of Go: methods](https://go.dev/tour/methods/1) (stop before interfaces)
- [Effective Go: getters](https://go.dev/doc/effective_go#Getters), [constructors](https://go.dev/doc/effective_go#composite_literals)
- [Go Code Review Comments: receiver type](https://go.dev/wiki/CodeReviewComments#receiver-type)
