package main

import (
	"fmt"
	"math"
)

func main() {
	var price, quantity int = 100, 4;	
	fmt.Printf("Price: %d, Quantity: %d\n", price, quantity);

	var (
		name string = "Laptop"
		q int = 5
	)

	fmt.Printf("Name: %s, q-Quantity: %d\n", name, q);

	a,b := 10, "string"

	fmt.Printf("a: %d, b: %s\n", a, b);

	x,y := 156.36,100.25

	z := math.Max(x,y)

	fmt.Printf("Max: %f between %f and %f\n", z,x,y);

}