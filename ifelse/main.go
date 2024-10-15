package main

import "fmt"

func main() {
	ticketPrice := 0;

	if age := 10; age <= 5 {
		ticketPrice = 0
	}else if age >= 5 && age < 22{
		ticketPrice = 10
	}else{
		ticketPrice = 15
	}

	fmt.Printf("%d", ticketPrice)
}