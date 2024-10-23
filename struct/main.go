package main

import (
	"fmt"
	"struct/computer"
)

type Employee struct {
	firstName string
	lastName  string
	age       int
	address	Address
}

type Address struct {
	city string 
	state string
	zip int
}

// methods on struct

type Circle struct {
	radius float64
}

func (c Circle) area() float64 { 
	a := 3.14 * c.radius * c.radius
	return a                      
}

func main() {
	e1 := Employee{
		firstName: "Sam",
		lastName:  "Anderson",
		age:       55,
		address:  Address{
			city: "New York",
			state: "NY",
			zip: 12345,
		},
	}


	fmt.Printf("Employee 1: %v\n", e1)

	c1 := computer.Computer{
		Brand: "Apple",
		Price: 2000,
	}

	fmt.Printf("Computer 1: %v\n", c1)

	// Annonymous struct											

	anonymusE := struct{
		firstName string
		middleName string
		lastName string
		age int
	}{
		firstName: "Sam",
		middleName: "John",
		lastName: "Anderson",
		age: 55,
	}

	fmt.Printf("Annonymous Employee: %v\n", anonymusE)
	fmt.Println(anonymusE.lastName)

	var addressOfE *Employee = &e1

	fmt.Printf("Address of e1: %v\n", addressOfE)
	fmt.Printf("Value of e1: %v\n", *addressOfE)

	// Equlity of two structs

	e2 := Employee{
		firstName: "Sam",
		lastName:  "Anderson",
		age:       55,
	}
	e2_dup := Employee{
		firstName: "Sam",
		lastName:  "Anderson",
		age:       55,
	}

	// a1 := Address{
	// 	city: "New York",
	// 	state: "NY",
	// 	zip: 12345,
	// }	

	if(e2 == e2_dup){
		fmt.Println("e1 and e2 are equal") // all the fields of the struct should be equal to get this output
	}

	// if (a1 == e1) cannot compare two different structs

	// Methos on struct
	cir1 := Circle{
		radius: 5,
	}
	fmt.Printf("Area of circle: %v\n", cir1.area())

	// built func from methods
	areaF := func (c Circle) float64 { 
		a := 3.14 * c.radius * c.radius
		return a                      
	}

	fmt.Printf("Area of circle: %v\n", areaF(cir1))

	
	

}