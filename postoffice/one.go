package main

import (
	"fmt"
	"time"
)

type Customer struct {
	Name string
	Item string
}

type Worker struct{
	Name string
}

func (c *Customer) Giveaway(){
	fmt.Printf("customer: %s has given away the parcel: %s\n",c.Name,c.Item)
}

func (w *Worker) Receive(item string){
	fmt.Printf("worker: %s has 'received' the parcel containing: %s\n",w.Name, item)
	fmt.Printf("worker: %s has 'started processing' the parcel containing: %s\n",w.Name, item)

	time.Sleep(2 * time.Second)

	fmt.Printf("worker: %s has 'processed' the parcel containing: %s\n",w.Name, item)
}

func main() {
	bobWorker := Worker{Name: "Bob"}

	zlatan := Customer{Name: "Zlatan", Item: "football"}
	ben := Customer{Name: "Ben", Item: "box"}
	jenny := Customer{Name: "Jenny", Item: "watermelon"}
	eric := Customer{Name: "Eric", Item: "teddy bear"}
	lisa := Customer{Name: "Lisa", Item: "basketball"}

	q := []Customer{zlatan, ben, jenny, eric, lisa}

	for _, c := range q {
		c.Giveaway()
		bobWorker.Receive(c.Item)
	}

}