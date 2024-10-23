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

type DeskWorker struct{
	Name string
	BackOfficeDeskChan chan string
}
func (w *DeskWorker) startDayShift(deskChan chan string, phoneChan chan string, shutdownChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
			case item,ok := <- deskChan:
				if ok {
					time.Sleep(1 * time.Second) 
					w.BackOfficeDeskChan <- item
					fmt.Printf("Desk worker %s started passing %s to the back office...\n", w.Name, item)
					fmt.Printf("Desk worker %s passed %s to the back office\n\n", w.Name, item)
				}
			case call,ok:= <- phoneChan: 
				if ok {
					fmt.Printf("CALL --->  %s\n", call)
				}
			case <-shutdownChan:
				fmt.Printf("the desk is closed - time to go home, %s\n", w.Name)
				return
		}
	}
	
	
	fmt.Println("the desk is closed - time to go home")
}
func (w *DeskWorker) Receive(item string){
	fmt.Printf("Deskworker: %s has 'received' the parcel containing: %s\n",w.Name, item)
	fmt.Printf("Deskworker: %s has 'started checking ID of the customer' the parcel containing: %s\n",w.Name, item)
}

type BackOfficeWorker struct {
 	Name string
}
func (bow *BackOfficeWorker) startDayShift(backOfficeDeskChan chan string, shutdownChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
			case item := <-backOfficeDeskChan:
				bow.Process(item)
			case <-shutdownChan:
				fmt.Printf("the back office is closed - time to go home, %s\n", bow.Name)
			return
		}
	}
}
func (bow *BackOfficeWorker) Process(item string) {
 fmt.Printf("Back office worker %s received %s\n", bow.Name, item)
 fmt.Printf("Back office worker %s started processing %s...\n", bow.Name, item)

 // to simulate long processing
 time.Sleep(10 * time.Second)

 fmt.Printf("Back office worker %s finished processing %s\n", bow.Name, item)
}




func main() {
	start := time.Now()

	wg := sync.WaitGroup{}

	deskChan := make(chan string, 3)
	backOfficeDeskChan := make(chan string, 3)
	phoneChan := make(chan string, 1) 

	deskShutdownChan := make(chan struct{})
	backOfficeDeskShutdownChan := make(chan struct{})

	bobWorker := DeskWorker{Name: "Bob", BackOfficeDeskChan: backOfficeDeskChan}
	odaWorker := BackOfficeWorker{Name: "Oda"}
	robertWorker := BackOfficeWorker{Name: "Robert"}
	marthaWorker := BackOfficeWorker{Name: "Martha"}

	zlatan := Customer{Name: "Zlatan", Item: "football"}
	ben := Customer{Name: "Ben", Item: "box"}
	jenny := Customer{Name: "Jenny", Item: "watermelon"}
	eric := Customer{Name: "Eric", Item: "teddy bear"}
	lisa := Customer{Name: "Lisa", Item: "basketball"}

	q := []Customer{zlatan, ben, jenny, eric, lisa}

	go func() {
		phoneChan <- "package arrived yet?"
		time.Sleep(1 * time.Second)
		phoneChan <- "now?"
	}()
	
	// DESK WORKERS

	go func () {
		wg.Add(1)
		bobWorker.startDayShift(deskChan, phoneChan, deskShutdownChan, &wg)
	}()

	// ----------------------------- //

	// BACK OFFICE WORKERS

	go func () {
		wg.Add(1)
		odaWorker.startDayShift(backOfficeDeskChan, backOfficeDeskShutdownChan, &wg)
	}()
	go func () {
		wg.Add(1)
		robertWorker.startDayShift(backOfficeDeskChan, backOfficeDeskShutdownChan, &wg)
	}()
	go func () {
		wg.Add(1)
		marthaWorker.startDayShift(backOfficeDeskChan, backOfficeDeskShutdownChan, &wg)
	}()

	// ----------------------------- //

	for _, c := range q {
		bobWorker.Receive(c.Item)
		deskChan <- c.Item
		fmt.Printf("customer: %s has 'queued' the parcel: %s\n",c.Name,c.Item)
	}

	deskShutdownChan <- struct{}{}

	for i := 0; i < 3; i++ {
		backOfficeDeskShutdownChan <- struct{}{}
	}

	close(backOfficeDeskChan)
	close(deskShutdownChan)
	close(deskChan)
	close(phoneChan)

	wg.Wait()

	end := time.Now()

	fmt.Println("Time taken: ", end.Sub(start).Seconds())
}