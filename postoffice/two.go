package main

import (
	"fmt"
	"time"
	"sync"
)

type Customer struct {
	Name string
	Item string
}

type Worker struct{
	Name string
	isDoneForTheDay bool
}

func (w *Worker) Receive(item string){
	fmt.Printf("worker: %s has 'received' the parcel containing: %s\n",w.Name, item)
	fmt.Printf("worker: %s has 'started processing' the parcel containing: %s\n",w.Name, item)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("worker: %s has 'processed' the parcel containing: %s\n",w.Name, item)
}

func (w *Worker) WaitToFinish() {
 for !w.isDoneForTheDay {}

 fmt.Printf("Worker %s has finished work for today\n", w.Name)
}

func (w *Worker) startDayShift(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	w.isDoneForTheDay = false
	for {
		  item, ok := <-ch
			if !ok {
				break
			}
		w.Receive(item)
	}
	w.isDoneForTheDay = true
	fmt.Println("the desk is closed - time to go home")
}


func main() {
	wg := sync.WaitGroup{}

	bobWorker := Worker{Name: "Bob"}

	zlatan := Customer{Name: "Zlatan", Item: "football"}
	ben := Customer{Name: "Ben", Item: "box"}
	jenny := Customer{Name: "Jenny", Item: "watermelon"}
	eric := Customer{Name: "Eric", Item: "teddy bear"}
	lisa := Customer{Name: "Lisa", Item: "basketball"}

	q := []Customer{zlatan, ben, jenny, eric, lisa}

	deskChan := make(chan string, 1)
	
	go func () {
		wg.Add(1)
		bobWorker.startDayShift(deskChan, &wg)
	}()

	for _, c := range q {
		deskChan <- c.Item
		fmt.Printf("customer: %s has 'queued' the parcel: %s\n",c.Name,c.Item)
	}
	
	close(deskChan)

	wg.Wait()
}