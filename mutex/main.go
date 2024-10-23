package main

import (
	"fmt"
	"time"
	"sync"
)

type Balance struct {
	amount int 
	currency string
	mu sync.Mutex
}
func (b *Balance) Write(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.amount += amount
}
func (b *Balance) Read() int {
	return b.amount
}

func main() {
	b := Balance{amount: 0, currency: "INR", mu: sync.Mutex{}}

	for i := 0; i < 1000; i++ {
		go func() {
			b.Write(1)
		}()
	}

	time.Sleep(2 * time.Second)
	fmt.Println(b.Read())

}