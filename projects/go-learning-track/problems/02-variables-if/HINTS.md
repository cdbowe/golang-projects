# Hints - 02

<details>
<summary>Hint 1 - direction</summary>

`Tier` is a range check: reject the invalid scores first, then walk the bands in one direction so each
test only needs a single comparison. `RiskBand` is a lookup on a `string` — the README says use a
`switch`, and several tiers share a band.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- A conditionless `switch` (`switch { case score >= 800: ... }`) replaces an `if`/`else if` chain. See
  the LESSON "Syntax" block.
- One `case` can list several values: `case "Fair", "Good":`.
- Don't forget `default` in `RiskBand` — the tests pass `"Platinum"` and `""`.
- `Tier(0)` must be `Invalid`, so your guard has to run before the band checks.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
func Tier(score int) string:
    if score < 300 or score > 850: return Invalid
    switch:
      case score >= 800: Exceptional
      case score >= 740: Very Good
      case score >= 670: Good
      case score >= 580: Fair
      default:           Poor      // 300..579 is all that's left

func RiskBand(tier string) string:
    switch tier:
      case Poor, Invalid:            Decline
      case Fair, Good:               Manual review
      case Very Good, Exceptional:   Auto approve
      default:                       Unknown
```

The exact strings are in README.md's tables — copy them character for character.

</details>
