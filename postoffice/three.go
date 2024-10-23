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
}

func (w *Worker) Receive(item string){
	fmt.Printf("worker: %s has 'received' the parcel containing: %s\n",w.Name, item)
	fmt.Printf("worker: %s has 'started processing' the parcel containing: %s\n",w.Name, item)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker: %s has 'processed' the parcel containing: %s\n",w.Name, item)
}

func (w *Worker) startDayShift(deskChan chan string, phoneChan chan string,wg *sync.WaitGroup) {
	defer wg.Done()
	keepWorking := true

	for keepWorking{
		select {
			case item,ok := <- deskChan:
				if ok {
					w.Receive(item)
				}else{
					keepWorking = false
				}
			case call,ok:= <- phoneChan: 
				if ok {
					fmt.Printf("CALL --->  %s\n", call)
				}else{
						keepWorking = false
				}
		}
	}
	
	
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
	phoneChan := make(chan string, 1)

	go func() {
		phoneChan <- "package arrived yet?"
		time.Sleep(1 * time.Second)
		phoneChan <- "now?"
	}()
	
	go func () {
		wg.Add(1)
		bobWorker.startDayShift(deskChan,phoneChan, &wg)
	}()

	for _, c := range q {
		deskChan <- c.Item
		fmt.Printf("customer: %s has 'queued' the parcel: %s\n",c.Name,c.Item)
	}


	close(deskChan)
	close(phoneChan)

	wg.Wait()
}