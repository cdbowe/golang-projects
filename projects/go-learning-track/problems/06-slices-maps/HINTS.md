# Hints - 06

<details>
<summary>Hint 1 - direction</summary>

- `GroupByZip`: guard the lengths, make the map, then one `range` over `zips` using the index to reach
  the matching price.
- `AverageFor`: comma-ok lookup, then guard for "absent or empty" in one condition, then sum and divide.
- `SortedZips`: you can't sort a map. Collect the keys into a slice, then sort the slice.
- `Cheapest`: never sort `prices` itself. Work on a copy.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- `groups[zip] = append(groups[zip], prices[i])` works even the first time a ZIP appears — see Check
  yourself #2.
- `make(map[string][]float64)` for a non-nil map; `make([]string, 0, len(groups))` for the keys.
- `prices, ok := groups[zip]` then `if !ok || len(prices) == 0`.
- [`slices.Clone`](https://pkg.go.dev/slices#Clone), [`slices.Sort`](https://pkg.go.dev/slices#Sort),
  and the `min` built-in (Go 1.21+) for clamping `n`.
- Add `"slices"` to the import block in `main.go`.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
GroupByZip:
    if len(zips) != len(prices): return nil, ErrLengthMismatch
    groups := make(map...)
    for i, zip := range zips: groups[zip] = append(groups[zip], prices[i])
    return groups, nil

AverageFor:
    prices, ok := groups[zip]
    if !ok || len(prices) == 0: return 0, false
    sum the prices; return sum / float64(len(prices)), true

SortedZips:
    keys := make([]string, 0, len(groups))
    for k := range groups: keys = append(keys, k)
    slices.Sort(keys); return keys

Cheapest:
    if n <= 0: return an empty slice
    c := slices.Clone(prices); slices.Sort(c)
    return c[:min(n, len(c))]
```

</details>
