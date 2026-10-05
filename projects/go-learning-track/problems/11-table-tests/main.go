package main

import "fmt"

// todaysRates is a hard-coded RateSource for the demo in main.
type todaysRates struct{}

func (todaysRates) BaseRate(termMonths int) (float64, error) {
	switch termMonths {
	case 180:
		return 0.0575, nil
	case 360:
		return 0.0625, nil
	}
	return 0, fmt.Errorf("no published rate for %d months", termMonths)
}

func main() {
	applicants := []struct {
		credit int
		ltv    float64
		term   int
	}{
		{760, 0.75, 360},
		{700, 0.85, 180},
		{600, 0.80, 360},
	}

	for _, a := range applicants {
		rate, err := Quote(todaysRates{}, a.credit, a.ltv, a.term)
		if err != nil {
			fmt.Printf("credit %d, LTV %.2f, %d mo: rejected: %v\n", a.credit, a.ltv, a.term, err)
			continue
		}
		fmt.Printf("credit %d, LTV %.2f, %d mo: %.3f%%\n", a.credit, a.ltv, a.term, rate*100)
	}
}
