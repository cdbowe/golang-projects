package main

import "testing"

func TestGreetingText(t *testing.T) {
	got := Greeting()
	want := "Hello, World!"
	if got != want {
		t.Errorf("Greeting() = %q, expected %q", got, want)
	}
}

func TestGreetingNotEmpty(t *testing.T) {
	if Greeting() == "" {
		t.Error("Greeting() returned the empty string; it must return text")
	}
}

func TestGreetingWithName(t *testing.T) {
	var actual string = GreetingWithName("Foo")
	var expected string = "Hello, Foo!"

	if actual != expected {
		t.Errorf("GreetingWithName() = %q, expected %q", actual, expected)
	}
}

func TestGreetingWithNameNotEmpty(t *testing.T) {
	var actual string = GreetingWithName("Bar")
	if actual == "" {
		t.Errorf("GreetingWithName() returned empty string; it must return text")
	}
}
