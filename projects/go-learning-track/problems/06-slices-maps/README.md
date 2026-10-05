# 06 - Listings by ZIP

**New concept:** slices (`len`/`cap`/`append`), maps, `range` over both, the comma-ok idiom
**Builds on:** 03, 04, 05

## Task

Listings arrive as two parallel slices: `zips[i]` goes with `prices[i]`. Implement four functions.

| Function | Returns |
|---|---|
| `GroupByZip(zips, prices)` | `map[zip][]price`, prices in input order. Lengths differ → `(nil, ErrLengthMismatch)`. Success → never a nil map |
| `AverageFor(groups, zip)` | `(average, true)`. ZIP missing, or present with no prices → `(0, false)` |
| `SortedZips(groups)` | the map's keys, ascending |
| `Cheapest(prices, n)` | the `n` lowest prices, ascending. `n > len` → all. `n <= 0` → empty. **Must not modify or share memory with `prices`** |

`ErrLengthMismatch` is given. `slices`, `maps` and `errors` are all standard library.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | Prices grouped per ZIP, in input order | `TestGroupByZip` |
| 2 | Length mismatch → `ErrLengthMismatch` and a nil map | `TestGroupByZipLengthMismatch` |
| 3 | Empty input → empty, **non-nil** map | `TestGroupByZipEmpty` |
| 4 | Average for present ZIPs | `TestAverageFor` |
| 5 | Missing ZIP or nil map → `(0, false)` | `TestAverageForMissing` |
| 6 | Present-but-empty group → `(0, false)`, no division by zero | `TestAverageForEmptyGroup` |
| 7 | Keys sorted, every time, despite random map order | `TestSortedZips`, `TestSortedZipsEmpty` |
| 8 | `n` lowest, ascending, with clamping | `TestCheapest` |
| 9 | Caller's slice is unchanged afterwards | `TestCheapestLeavesInputAlone` |
| 10 | Result doesn't share a backing array with the input | `TestCheapestDoesNotShareMemory` |

## Run

```bash
go test ./problems/06-slices-maps/...
go run ./problems/06-slices-maps
```

## Stretch (optional)

Rewrite `SortedZips` as one line using `slices.Sorted` and `maps.Keys`. Same tests must pass.
