# 02 - Credit score tiers

**New concept:** declarations (`var` vs `:=`), zero values, `if`/`else if`, `switch`
**Reference:** `LESSON.md` also has the full `fmt` format-verb table
**Builds on:** 01

## Task

Implement two functions.

`Tier(score int) string` maps a score to a tier:

| Score | Tier |
|---|---|
| below 300, or above 850 | `Invalid` |
| 300-579 | `Poor` |
| 580-669 | `Fair` |
| 670-739 | `Good` |
| 740-799 | `Very Good` |
| 800-850 | `Exceptional` |

`RiskBand(tier string) string` maps a tier name to a decision band with a `switch`:

| Tier | Band |
|---|---|
| `Poor`, `Invalid` | `Decline` |
| `Fair`, `Good` | `Manual review` |
| `Very Good`, `Exceptional` | `Auto approve` |
| anything else | `Unknown` |

Bounds are inclusive. `RiskBand` takes the *tier name*, not the score.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | Every tier boundary (300/579/580/669/670/739/740/799/800/850) returns the right tier | `TestTierBoundaries` |
| 2 | Out-of-range scores, including `0` and negatives, return `Invalid` | `TestTierInvalid` |
| 3 | `RiskBand` maps all six known tiers, and unknown input to `Unknown` | `TestRiskBand` |
| 4 | `RiskBand(Tier(712))` is `Manual review` | `TestPipeline` |

## Run

```bash
go test ./problems/02-variables-if/...
go run ./problems/02-variables-if
```

## Stretch (optional)

Group two cases on one `switch` line (`case "Poor", "Invalid":`) if you wrote them separately — then
confirm the tests still pass.
