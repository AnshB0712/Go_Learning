package main

import "fmt"

type myString string

func main() {
	const a int = 10

	const (
		LIMIT = 10
		AUTH  = "none"
	)

	fmt.Printf("a: %d, LIMIT: %d, AUTH: %s\n", a, LIMIT, AUTH)

	// value constants should be known at compile time not in runtime
	// const b int = math.Sqrt(4) --> Wrong

	// something about untyped constants to remember about
	const hello = "hello world"
	// const n = hello
	fmt.Printf("hello: %s, typeof: %T\n", hello, hello)

	// ------------------------------------------------ //

	var defaultName = "Sam"

	var someName myString = "Sam"
	fmt.Printf("defaultName: %s, someName: %s\n", defaultName, someName)
	fmt.Printf("defaultName type: %T, someName type: %T\n", defaultName, someName)

	// cannot do this as someName is of type myString and other is string --> custom type and default type are different in go nature eventhough underlying type is same
	// someName = defaultName

	// ------------------------------------------------ //
	
	const c = 5
	var intVar int = c
	var int32Var int32 = c
	var float64Var float64 = c

	fmt.Printf("intVar: %d, int32Var: %d, float64Var: %f\n", intVar, int32Var, float64Var)


}