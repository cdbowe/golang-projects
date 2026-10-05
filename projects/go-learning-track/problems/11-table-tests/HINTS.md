# Hints - 11

<details>
<summary>Hint 1 - direction</summary>

Turn `Quote`'s doc comment into a list of behaviours, then into rows. Every comparison in the code
(`< 620`, `>= 740`, `>= 680`, `> 0.95`, `> 0.80`, the term check) gets a row on **each side** of its
boundary. Then add rows for each error path, including the source failing.

Two fakes cover everything: one that answers per term, and one that always fails.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- Row fields: `name string`, `src RateSource`, `credit int`, `ltv float64`, `term int`, `want float64`,
  `wantErr error`.
- Per-term fake: `type termRates map[int]float64` with a `BaseRate` method that looks the term up
  (comma-ok, 06) and returns an error for unknown terms. Give 180 and 360 **different** rates.
- Failing fake: a struct holding an `err error` field that `BaseRate` returns. Use a package-level
  `var errSourceDown = errors.New("source down")` so you can check it with `errors.Is`.
- `func assertRate(t *testing.T, got, want float64)` with `t.Helper()` and `math.Abs(got-want) > 1e-9`.
- Exact boundaries to include: credit 619/620, 679/680, 739/740; LTV 0.80/0.81 and 0.95/0.96; term 180,
  360, and something else.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
var errSourceDown = errors.New(...)
type termRates map[int]float64          // + BaseRate method
type failingSource struct{ err error }  // + BaseRate method

func assertRate(t, got, want): t.Helper(); compare with tolerance

func TestQuote(t):
    rates := termRates{180: 0.05, 360: 0.06}
    tests := []struct{ name; src; credit; ltv; term; want; wantErr }{
        {"top tier, low LTV, 30y",  rates, 740, 0.80, 360, 0.06,   nil},
        {"top tier, 15y",           rates, 740, 0.80, 180, 0.05,   nil},
        {"just below top tier",     rates, 739, 0.80, 360, 0.0625, nil},
        ... one row per boundary side ...
        {"credit below floor",      rates, 619, 0.80, 360, 0,      ErrIneligible},
        {"unsupported term",        rates, 760, 0.80, 240, 0,      ErrUnsupportedTerm},
        {"source down",             failingSource{errSourceDown}, 760, 0.80, 360, 0, ErrRateUnavailable},
    }
    for each tc: t.Run(tc.name, func(t):
        got, err := Quote(tc.src, tc.credit, tc.ltv, tc.term)
        if tc.wantErr != nil: assert errors.Is(err, tc.wantErr); return
        assert err == nil; assertRate(t, got, tc.want)
    )

// plus: a separate check that the "source down" error also Is errSourceDown
```

</details>
