package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// bool
	a := true
	b := false

	fmt.Printf("a: %t, b: %t\n", a, b)

	c := !a || b || false;

	fmt.Printf("c: %t\n", c)

	d := a && b

	fmt.Printf("d: %t\n", d)

	// int

	e := 128

	fmt.Printf("e is of type %T with value od %d with size of %d bytes\n", e,e,unsafe.Sizeof(e))

	f := e * e * e

	fmt.Printf("f is of type %T with value od %d with size of %d bytes\n", f,f,unsafe.Sizeof(f))

	// string

	contender := "Max Holloway"
	g := "UF"
	h := "C"
	event := g+h

	fmt.Printf("contender with name: %s will fight in %s 308\n", contender, event)

	// type conversion
	intNum := 10
	floatNum := 10.525

	// cannot do it as its different types of data
	// sum := intNum + floatNum
	// instead convert type to same
	sum := float64(intNum) + floatNum
	sum2 := intNum + int(floatNum)


	fmt.Printf("sum is of type %T with value of %f\n", sum, sum)
	fmt.Printf("sum2 is of type %T with value of %d\n", sum2, sum2)
}