package main

import (
	"fmt"
	"time"
	"sync"
)

type Semaphore struct {
	ch chan int
}

func (s *Semaphore) Aquire() {
	s.ch <- 1
}
func (s *Semaphore) Release() {
	<- s.ch
}

func createSemaphore(limit int) *Semaphore {

	if limit < 0 {
		limit = 0
	}

	s := Semaphore{
		ch: make(chan int, limit),
	}

	return &s
}

func main() {
	var wg sync.WaitGroup

	semaphore := createSemaphore(10)

	for v := range 100 {
		wg.Add(1) 
		go func(){
			defer wg.Done()
			semaphore.Aquire()
			fmt.Printf("Request: %d at Time: %v \n", v, time.Now())
			time.Sleep(2 * time.Second)
			semaphore.Release()
		}()
	}

	wg.Wait()
}