package main

import (
	"errors"
	"math"
	"slices"
	"testing"
)

func closeEnough(got, want float64) bool {
	return math.Abs(got-want) < 0.01
}

func TestGroupByZip(t *testing.T) {
	groups, err := GroupByZip(
		[]string{"30301", "30305", "30301"},
		[]float64{300000, 450000, 350000},
	)
	if err != nil {
		t.Fatalf("GroupByZip returned error %v, want nil", err)
	}
	if len(groups) != 2 {
		t.Errorf("len(groups) = %d, want 2", len(groups))
	}
	if got, want := groups["30301"], []float64{300000, 350000}; !slices.Equal(got, want) {
		t.Errorf("groups[%q] = %v, want %v (input order)", "30301", got, want)
	}
	if got, want := groups["30305"], []float64{450000}; !slices.Equal(got, want) {
		t.Errorf("groups[%q] = %v, want %v", "30305", got, want)
	}
}

func TestGroupByZipLengthMismatch(t *testing.T) {
	groups, err := GroupByZip([]string{"30301"}, []float64{300000, 450000})
	if !errors.Is(err, ErrLengthMismatch) {
		t.Errorf("err = %v, want ErrLengthMismatch", err)
	}
	if groups != nil {
		t.Errorf("groups = %v, want nil on error", groups)
	}
}

func TestGroupByZipEmpty(t *testing.T) {
	groups, err := GroupByZip(nil, nil)
	if err != nil {
		t.Fatalf("GroupByZip(nil, nil) returned error %v, want nil", err)
	}
	if groups == nil {
		t.Fatal("groups is a nil map; on success return an empty, non-nil map (writing to a nil map panics)")
	}
	if len(groups) != 0 {
		t.Errorf("len(groups) = %d, want 0", len(groups))
	}
}

func TestAverageFor(t *testing.T) {
	groups := map[string][]float64{
		"30301": {300000, 350000},
		"30305": {450000},
	}

	avg, ok := AverageFor(groups, "30301")
	if !ok || !closeEnough(avg, 325000) {
		t.Errorf("AverageFor(%q) = (%v, %v), want (325000, true)", "30301", avg, ok)
	}

	avg, ok = AverageFor(groups, "30305")
	if !ok || !closeEnough(avg, 450000) {
		t.Errorf("AverageFor(%q) = (%v, %v), want (450000, true)", "30305", avg, ok)
	}
}

func TestAverageForMissing(t *testing.T) {
	groups := map[string][]float64{"30301": {300000}}

	if avg, ok := AverageFor(groups, "99999"); ok || avg != 0 {
		t.Errorf("AverageFor(missing zip) = (%v, %v), want (0, false)", avg, ok)
	}
	if avg, ok := AverageFor(nil, "30301"); ok || avg != 0 {
		t.Errorf("AverageFor(nil map) = (%v, %v), want (0, false)", avg, ok)
	}
}

// A ZIP that is present with no prices must not divide by zero.
func TestAverageForEmptyGroup(t *testing.T) {
	groups := map[string][]float64{
		"30309": {},
		"30310": nil,
	}

	if avg, ok := AverageFor(groups, "30309"); ok || avg != 0 {
		t.Errorf("AverageFor(empty slice) = (%v, %v), want (0, false)", avg, ok)
	}
	if avg, ok := AverageFor(groups, "30310"); ok || avg != 0 {
		t.Errorf("AverageFor(nil slice) = (%v, %v), want (0, false)", avg, ok)
	}
}

// Map iteration order is randomized, so this runs repeatedly: an unsorted
// result will not survive 20 attempts.
func TestSortedZips(t *testing.T) {
	groups := map[string][]float64{
		"30305": {1},
		"10001": {2},
		"30301": {3},
		"94110": {4},
	}
	want := []string{"10001", "30301", "30305", "94110"}

	for range 20 {
		if got := SortedZips(groups); !slices.Equal(got, want) {
			t.Fatalf("SortedZips() = %v, want %v", got, want)
		}
	}
}

func TestSortedZipsEmpty(t *testing.T) {
	if got := SortedZips(map[string][]float64{}); len(got) != 0 {
		t.Errorf("SortedZips(empty) = %v, want no keys", got)
	}
}

func TestCheapest(t *testing.T) {
	prices := []float64{500, 100, 300, 200}

	if got, want := Cheapest(prices, 2), []float64{100, 200}; !slices.Equal(got, want) {
		t.Errorf("Cheapest(prices, 2) = %v, want %v", got, want)
	}
	if got, want := Cheapest(prices, 10), []float64{100, 200, 300, 500}; !slices.Equal(got, want) {
		t.Errorf("Cheapest(prices, 10) = %v, want %v (n > len returns all)", got, want)
	}
	if got := Cheapest(prices, 0); len(got) != 0 {
		t.Errorf("Cheapest(prices, 0) = %v, want empty", got)
	}
	if got := Cheapest(prices, -1); len(got) != 0 {
		t.Errorf("Cheapest(prices, -1) = %v, want empty", got)
	}
}

func TestCheapestLeavesInputAlone(t *testing.T) {
	prices := []float64{500, 100, 300, 200}
	original := []float64{500, 100, 300, 200}

	_ = Cheapest(prices, 2)
	if !slices.Equal(prices, original) {
		t.Errorf("after Cheapest, prices = %v, want %v (sorting in place changed the caller's slice)", prices, original)
	}
}

func TestCheapestDoesNotShareMemory(t *testing.T) {
	prices := []float64{500, 100, 300, 200}

	got := Cheapest(prices, 4)
	if len(got) == 0 {
		t.Fatal("Cheapest(prices, 4) is empty")
	}
	got[0] = -1
	if slices.Contains(prices, -1) {
		t.Errorf("writing to the result changed prices = %v; the result shares a backing array with the input", prices)
	}
}
