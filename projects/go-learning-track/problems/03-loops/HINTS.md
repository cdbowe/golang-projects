# Hints - 03

<details>
<summary>Hint 1 - direction</summary>

`Label` is problem 02's `switch` again, over `month%3` and `month%4`. Order matters: the "divisible by
both" case has to be tested **before** the single-divisor cases, or 12 never reaches it.

`Schedule` is one loop calling `Label` and gluing lines together. Decide up front how you avoid the
trailing newline.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- Guard `month < 1` first, then `switch { case month%3 == 0 && month%4 == 0: ... }`.
- `fmt.Sprintf("%d: %s", month, Label(month))` builds one line as a string.
- For the loop, the C-style form starts where you need it: `for m := 1; m <= months; m++`.
  (`for m := range months` starts at 0 — you'd have to use `m+1` everywhere.)
- Joining: either `strings.Builder` with `WriteString`, or `+=` on a local string. Both pass. Write the
  separator *before* each line except the first.
- Adding `strings` to the import block: group it with `"fmt"` inside `import ( ... )`.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
func Label(month int) string:
    if month < 1: return Invalid
    switch:
      case month%3 == 0 && month%4 == 0: Escrow check + Statement audit
      case month%3 == 0:                 Escrow check
      case month%4 == 0:                 Statement audit
      default:                           Payment

func Schedule(months int) string:
    var b strings.Builder
    for m := 1; m <= months; m++:
        if m > 1: b.WriteString("\n")
        write "<m>: <Label(m)>" into b
    return b.String()
```

`months <= 0` needs no special case here — think about whether the loop body ever runs.

</details>
