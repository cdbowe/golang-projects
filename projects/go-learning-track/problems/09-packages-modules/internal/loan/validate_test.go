// White-box test: same package as the code, so unexported names are visible.
package loan

import "testing"

func TestValidID(t *testing.T) {
	if !validID("L-100") {
		t.Errorf("validID(%q) = false, want true", "L-100")
	}
	if validID("") {
		t.Errorf("validID(%q) = true, want false", "")
	}
	if validID(" \t ") {
		t.Errorf("validID(%q) = true, want false", " \t ")
	}
}

func TestValidAmount(t *testing.T) {
	if !validAmount(0.01) {
		t.Error("validAmount(0.01) = false, want true")
	}
	if validAmount(0) {
		t.Error("validAmount(0) = true, want false")
	}
	if validAmount(-1) {
		t.Error("validAmount(-1) = true, want false")
	}
}
