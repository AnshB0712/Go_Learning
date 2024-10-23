package main 

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan int, 5)

	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	ratelimit := time.Tick(2 * time.Second)

	for req := range requests {
		<-ratelimit
		fmt.Println("Request: ", req, time.Now())
	}

	burstyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func(){
		for range time.Tick(2 * time.Second){
			burstyLimiter <- time.Now()
		}
	}()

	requests2 := make(chan int, 5)

	for i:=0; i<5; i++ {
		requests2 <- i
	}
	close(requests2)

	for v := range requests2 {
		<- burstyLimiter
		fmt.Println("Request: ", v, time.Now())
	}

}