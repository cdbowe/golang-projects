package main

import "testing"

func TestLabelPlainMonths(t *testing.T) {
	if got := Label(1); got != "Payment" {
		t.Errorf("Label(1) = %q, want %q", got, "Payment")
	}
	if got := Label(7); got != "Payment" {
		t.Errorf("Label(7) = %q, want %q", got, "Payment")
	}
}

func TestLabelEscrow(t *testing.T) {
	if got := Label(3); got != "Escrow check" {
		t.Errorf("Label(3) = %q, want %q", got, "Escrow check")
	}
	if got := Label(9); got != "Escrow check" {
		t.Errorf("Label(9) = %q, want %q", got, "Escrow check")
	}
}

func TestLabelStatement(t *testing.T) {
	if got := Label(4); got != "Statement audit" {
		t.Errorf("Label(4) = %q, want %q", got, "Statement audit")
	}
	if got := Label(8); got != "Statement audit" {
		t.Errorf("Label(8) = %q, want %q", got, "Statement audit")
	}
}

func TestLabelBoth(t *testing.T) {
	want := "Escrow check + Statement audit"
	if got := Label(12); got != want {
		t.Errorf("Label(12) = %q, want %q", got, want)
	}
	if got := Label(24); got != want {
		t.Errorf("Label(24) = %q, want %q", got, want)
	}
}

func TestLabelInvalid(t *testing.T) {
	if got := Label(0); got != "Invalid" {
		t.Errorf("Label(0) = %q, want %q", got, "Invalid")
	}
	if got := Label(-3); got != "Invalid" {
		t.Errorf("Label(-3) = %q, want %q", got, "Invalid")
	}
}

func TestScheduleFirstYear(t *testing.T) {
	want := "1: Payment\n" +
		"2: Payment\n" +
		"3: Escrow check\n" +
		"4: Statement audit\n" +
		"5: Payment\n" +
		"6: Escrow check\n" +
		"7: Payment\n" +
		"8: Statement audit\n" +
		"9: Escrow check\n" +
		"10: Payment\n" +
		"11: Payment\n" +
		"12: Escrow check + Statement audit"

	if got := Schedule(12); got != want {
		t.Errorf("Schedule(12) =\n%s\n\nwant:\n%s", got, want)
	}
}

func TestScheduleOneMonth(t *testing.T) {
	want := "1: Payment"
	if got := Schedule(1); got != want {
		t.Errorf("Schedule(1) = %q, want %q", got, want)
	}
}

func TestScheduleNoTrailingNewline(t *testing.T) {
	got := Schedule(4)
	if len(got) == 0 {
		t.Fatal("Schedule(4) is empty")
	}
	if got[len(got)-1] == '\n' {
		t.Errorf("Schedule(4) ends with a newline: %q", got)
	}
}

func TestScheduleEmpty(t *testing.T) {
	if got := Schedule(0); got != "" {
		t.Errorf("Schedule(0) = %q, want %q", got, "")
	}
	if got := Schedule(-1); got != "" {
		t.Errorf("Schedule(-1) = %q, want %q", got, "")
	}
}

// Now the same tests on Schedule2

func TestSchedule2FirstYear(t *testing.T) {
	want := "1: Payment\n" +
		"2: Payment\n" +
		"3: Escrow check\n" +
		"4: Statement audit\n" +
		"5: Payment\n" +
		"6: Escrow check\n" +
		"7: Payment\n" +
		"8: Statement audit\n" +
		"9: Escrow check\n" +
		"10: Payment\n" +
		"11: Payment\n" +
		"12: Escrow check + Statement audit"

	if got := Schedule2(12); got != want {
		t.Errorf("Schedule2(12) =\n%s\n\nwant:\n%s", got, want)
	}
}

func TestSchedule2OneMonth(t *testing.T) {
	want := "1: Payment"
	if got := Schedule2(1); got != want {
		t.Errorf("Schedule2(1) = %q, want %q", got, want)
	}
}

func TestSchedule2NoTrailingNewline(t *testing.T) {
	got := Schedule2(4)
	if len(got) == 0 {
		t.Fatal("Schedule2(4) is empty")
	}
	if got[len(got)-1] == '\n' {
		t.Errorf("Schedule2(4) ends with a newline: %q", got)
	}
}

func TestSchedule2Empty(t *testing.T) {
	if got := Schedule2(0); got != "" {
		t.Errorf("Schedule2(0) = %q, want %q", got, "")
	}
	if got := Schedule2(-1); got != "" {
		t.Errorf("Schedule2(-1) = %q, want %q", got, "")
	}
}
