package main 

import (
	"fmt"
	"time"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	timer := time.NewTimer(5 * time.Second)
	iteration := 0
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-timer.C:
				fmt.Println("Timer expired")
				if iteration == 5 {
					fmt.Println("Process Finished")
					return
				}
				iteration++
				timer.Reset(5 * time.Second)
			default:
				fmt.Println("Iteration: ", iteration)
				time.Sleep(4 * time.Second)
			}
		}
	}()

	fmt.Println("Process Started")

	wg.Wait()


}