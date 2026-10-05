# 03 - Payment schedule labels

**New concept:** `for` — Go's only loop (C-style, while-style, `range`)
**Builds on:** 01, 02

## Task

A FizzBuzz in loan-servicing clothes.

`Label(month int) string` labels one month:

| Month | Label |
|---|---|
| divisible by 3 **and** 4 (12, 24, ...) | `Escrow check + Statement audit` |
| divisible by 3 | `Escrow check` |
| divisible by 4 | `Statement audit` |
| anything else | `Payment` |
| below 1 | `Invalid` |

`Schedule(months int) string` returns one line per month from 1 to `months`, formatted `N: Label`,
joined by `\n`, **no trailing newline**. `months <= 0` returns `""`.

```
Schedule(4) == "1: Payment\n2: Payment\n3: Escrow check\n4: Statement audit"
```

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | Plain months return `Payment` | `TestLabelPlainMonths` |
| 2 | Multiples of 3 and of 4 get their own labels | `TestLabelEscrow`, `TestLabelStatement` |
| 3 | Months divisible by both get the combined label | `TestLabelBoth` |
| 4 | `Label(0)` and negatives return `Invalid` | `TestLabelInvalid` |
| 5 | `Schedule(12)` matches the full expected block | `TestScheduleFirstYear` |
| 6 | No trailing newline; `Schedule(1)` is a single line | `TestScheduleNoTrailingNewline`, `TestScheduleOneMonth` |
| 7 | `Schedule(0)` and `Schedule(-1)` return `""` | `TestScheduleEmpty` |

## Run

```bash
go test ./problems/03-loops/...
go run ./problems/03-loops
```

## Stretch (optional)

Write `Schedule` a second way — if you built the string with `+=`, redo it with `strings.Builder`, or
the reverse. Same tests must pass.
