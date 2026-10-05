package main

import "fmt"

// Tier maps a FICO-style credit score to its tier name.
// Scores outside 300..850 are "Invalid".
func Tier(score int) string {
	switch {
	case score < 300 || score > 850:
		return "Invalid"
	case score >= 800:
		return "Exceptional"
	case score >= 740:
		return "Very Good"
	case score >= 670:
		return "Good"
	case score >= 580:
		return "Fair"
	default:
		return "Poor"
	}
}

// RiskBand maps a tier name to the decision band for that tier.
// Unknown tier names return "Unknown".
func RiskBand(tier string) string {
	switch tier {
	case "Poor", "Invalid":
		return "Decline"
	case "Fair", "Good":
		return "Manual review"
	case "Very Good", "Exceptional":
		return "Auto approve"
	default:
		return "Unknown"
	}
}

func main() {
	score := 712
	tier := Tier(score)
	fmt.Printf("score %d -> tier %q -> band %q\n", score, tier, RiskBand(tier))
}
