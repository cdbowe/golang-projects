# Hints - 01

<details>
<summary>Hint 1 - direction</summary>

The test compares `Greeting()` against a literal. Make the function hand back that same literal
instead of `""`. No `fmt` needed inside `Greeting` — just a `return`.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

Read the `want` value in `TestGreetingText` in `main_test.go`. It is exact: capital H, capital W,
comma, space, exclamation mark. String literals in Go use double quotes (`"..."`); single quotes are
runes, not strings.

For `main()`, the function you want is [`fmt.Println`](https://pkg.go.dev/fmt#Println), already wired
up for you.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
func Greeting() string:
    return <the exact string the test wants>
```

Then run `go test ./problems/01-hello-world/...` and expect `ok`. Then `go run ./problems/01-hello-world`.

</details>
