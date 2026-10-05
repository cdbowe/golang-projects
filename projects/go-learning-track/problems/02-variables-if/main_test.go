package main

import "testing"

func TestTierBoundaries(t *testing.T) {
	if got := Tier(300); got != "Poor" {
		t.Errorf("Tier(300) = %q, want %q", got, "Poor")
	}
	if got := Tier(579); got != "Poor" {
		t.Errorf("Tier(579) = %q, want %q", got, "Poor")
	}
	if got := Tier(580); got != "Fair" {
		t.Errorf("Tier(580) = %q, want %q", got, "Fair")
	}
	if got := Tier(669); got != "Fair" {
		t.Errorf("Tier(669) = %q, want %q", got, "Fair")
	}
	if got := Tier(670); got != "Good" {
		t.Errorf("Tier(670) = %q, want %q", got, "Good")
	}
	if got := Tier(739); got != "Good" {
		t.Errorf("Tier(739) = %q, want %q", got, "Good")
	}
	if got := Tier(740); got != "Very Good" {
		t.Errorf("Tier(740) = %q, want %q", got, "Very Good")
	}
	if got := Tier(799); got != "Very Good" {
		t.Errorf("Tier(799) = %q, want %q", got, "Very Good")
	}
	if got := Tier(800); got != "Exceptional" {
		t.Errorf("Tier(800) = %q, want %q", got, "Exceptional")
	}
	if got := Tier(850); got != "Exceptional" {
		t.Errorf("Tier(850) = %q, want %q", got, "Exceptional")
	}
}

func TestTierInvalid(t *testing.T) {
	if got := Tier(299); got != "Invalid" {
		t.Errorf("Tier(299) = %q, want %q", got, "Invalid")
	}
	if got := Tier(851); got != "Invalid" {
		t.Errorf("Tier(851) = %q, want %q", got, "Invalid")
	}
	if got := Tier(0); got != "Invalid" {
		t.Errorf("Tier(0) = %q, want %q (the int zero value is not a valid score)", got, "Invalid")
	}
	if got := Tier(-50); got != "Invalid" {
		t.Errorf("Tier(-50) = %q, want %q", got, "Invalid")
	}
}

func TestRiskBand(t *testing.T) {
	if got := RiskBand("Poor"); got != "Decline" {
		t.Errorf("RiskBand(%q) = %q, want %q", "Poor", got, "Decline")
	}
	if got := RiskBand("Fair"); got != "Manual review" {
		t.Errorf("RiskBand(%q) = %q, want %q", "Fair", got, "Manual review")
	}
	if got := RiskBand("Good"); got != "Manual review" {
		t.Errorf("RiskBand(%q) = %q, want %q", "Good", got, "Manual review")
	}
	if got := RiskBand("Very Good"); got != "Auto approve" {
		t.Errorf("RiskBand(%q) = %q, want %q", "Very Good", got, "Auto approve")
	}
	if got := RiskBand("Exceptional"); got != "Auto approve" {
		t.Errorf("RiskBand(%q) = %q, want %q", "Exceptional", got, "Auto approve")
	}
	if got := RiskBand("Invalid"); got != "Decline" {
		t.Errorf("RiskBand(%q) = %q, want %q", "Invalid", got, "Decline")
	}
	if got := RiskBand("Platinum"); got != "Unknown" {
		t.Errorf("RiskBand(%q) = %q, want %q", "Platinum", got, "Unknown")
	}
	if got := RiskBand(""); got != "Unknown" {
		t.Errorf("RiskBand(%q) = %q, want %q", "", got, "Unknown")
	}
}

// TestPipeline is the end-to-end path main() uses.
func TestPipeline(t *testing.T) {
	if got := RiskBand(Tier(712)); got != "Manual review" {
		t.Errorf("RiskBand(Tier(712)) = %q, want %q", got, "Manual review")
	}
}
