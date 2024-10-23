package main

import (
	"fmt"
	"math/rand"
)

func main() {

	finger := 6
	switch finger {
	case 1:
		fmt.Println("Thumb")
	case 2: 
		fmt.Println("Index")
	case 3:
		fmt.Println("Middle")
	case 4:
		fmt.Println("Ring")
	case 5:
		fmt.Println("Pinky")
	default:
		fmt.Println("Human?")
	}

	letter := "i"
	switch letter {
	case "a","i","o","e","u":
		fmt.Println("vowel")
	default: 
		fmt.Println("not a vowel")
	}

	switch n := 7; {
	case n < 8: 
		fmt.Println("Less than 8")
		fallthrough
	case n < 9:
		fmt.Println("Less than 9")
	case n < 10:
		fmt.Println("less than 10")
		fallthrough
	default: 
		fmt.Println("done")
	}

	// fallthrough will go to the next case-block-function even though it evaluates to false
	switch num:=25; {
	case num < 50:
		fmt.Println("less than 50")
		fallthrough
	case num > 100:
		fmt.Println("greater than 100")
	}

	// switch inside loop common pattern that we see in go.
	loopylooop:
		for {
			switch n:=rand.Intn(10); {
			case n%2==0:
					break loopylooop;
			default: 
			fmt.Printf("%d\n", n)
			}
		}
}