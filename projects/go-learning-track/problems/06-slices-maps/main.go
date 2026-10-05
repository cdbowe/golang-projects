package main

import (
	"errors"
	"fmt"
	"maps"
	"slices"
)

// ErrLengthMismatch means the two input slices don't line up one-to-one.
var ErrLengthMismatch = errors.New("zips and prices have different lengths")

// GroupByZip pairs zips[i] with prices[i] and groups the prices by ZIP code,
// keeping each ZIP's prices in input order. It returns ErrLengthMismatch (and a
// nil map) when the slices differ in length. On success the map is never nil.
func GroupByZip(zips []string, prices []float64) (map[string][]float64, error) {
	if len(zips) != len(prices) {
		return nil, ErrLengthMismatch
	}

	zipToPrice := make(map[string][]float64, len(zips))

	for i, zip := range zips {
		zipToPrice[zip] = append(zipToPrice[zip], prices[i])
	}

	return zipToPrice, nil
}

// AverageFor returns the average price for zip. ok is false when zip is not in
// groups or has no prices.
func AverageFor(groups map[string][]float64, zip string) (avg float64, ok bool) {
	zipPrices, ok := groups[zip]

	if !ok || len(zipPrices) == 0 {
		return 0, false
	}

	sum := 0.0
	for _, v := range zipPrices {
		sum += v
	}

	return sum / float64(len(zipPrices)), true
}

// SortedZips returns the keys of groups in ascending order.
func SortedZips(groups map[string][]float64) []string {
	return slices.Sorted(maps.Keys(groups))

	// zips := make([]string, 0, len(groups))
	// for zip := range groups {
	// 	zips = append(zips, zip)
	// }

	// slices.Sort(zips)

	// return zips
}

// Cheapest returns the n lowest prices in ascending order, without modifying
// prices. n larger than len(prices) returns all of them; n <= 0 returns an
// empty slice.
func Cheapest(prices []float64, n int) []float64 {
	if n <= 0 {
		return []float64{}
	}

	pricesCopy := slices.Sorted(slices.Values(prices)) // Not to be confused with maps.Values()
	// pricesCopy := make([]float64, 0, len(prices))
	// for _, price := range prices {
	// 	pricesCopy = append(pricesCopy, price)
	// }

	// slices.Sort(pricesCopy)

	numPrices := min(n, len(pricesCopy))
	return pricesCopy[:numPrices]
}

func main() {
	zips := []string{"30301", "30305", "30301", "10001", "30305"}
	prices := []float64{310000, 455000, 289000, 1250000, 470000}

	groups, err := GroupByZip(zips, prices)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, zip := range SortedZips(groups) {
		avg, _ := AverageFor(groups, zip)
		fmt.Printf("%s  avg $%.0f  (%d listings)\n", zip, avg, len(groups[zip]))
	}
	fmt.Println("3 cheapest:", Cheapest(prices, 3))
	fmt.Println("input untouched:", prices)
}
