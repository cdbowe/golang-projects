package main

import "fmt"

// Greeting returns the text that main prints.
func Greeting() string {
	return "Hello, World!"
}

// GreetingWithName returns the text that main prints.
func GreetingWithName(name string) string {
	return "Hello, " + name + "!"
}

func main() {
	fmt.Println(Greeting())
	fmt.Println(GreetingWithName("Chris"))
}
