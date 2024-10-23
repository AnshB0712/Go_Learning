package main

import "fmt"

func main() {
	var a int = 10

	var addressOfa *int = &a

	fmt.Printf("Address of a: %v\n", addressOfa)
	fmt.Printf("Value of a: %v\n", *addressOfa)
}