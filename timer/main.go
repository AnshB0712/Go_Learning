package main 

import (
	"fmt"
	"time"
)

func main() {
	timer1 := time.NewTimer(5 * time.Second)

	go func() {
		fmt.Println("out of main thread")
	}()

	fmt.Println("Timer 1 started")
	<-timer1.C // Blocks until the timer expires
	fmt.Println("Timer 1 expired")

	fmt.Println("ansh")
}