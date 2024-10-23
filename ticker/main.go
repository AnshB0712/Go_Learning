package main 

import (
	"fmt"
	"time"
	// "sync"
)

// func main() {
// 	var wg sync.WaitGroup
// 	ticker := time.NewTicker(5 * time.Second)
// 	iteration := 0

// 	wg.Add(1)

// 	go func() {
// 		defer wg.Done()
// 		for {
// 			select {
// 				case <-ticker.C:
// 					fmt.Println("Ticker duration hit")
// 					if iteration == 5 {
// 						fmt.Println("Process Finished")
// 						ticker.Stop()
// 						return
// 					}
// 					iteration++
// 				default:	
// 					fmt.Println("Iteration: ", iteration)
// 					time.Sleep(5 * time.Second)
// 			}
// 		}
// 	}()
	
	

// 	fmt.Println("Process Started")

// 	wg.Wait()
// }

func main() {
	tick := time.Tick(5 * time.Second)

for v := range tick {
	fmt.Println("Tick at: ", v)
}

}