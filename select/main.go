package main 

import (
	"fmt" 
	"time"
)

func markChannelDone(ch chan <- string){
	time.Sleep(10 * time.Second)
	ch <- "Done"
}

func main() {
	ch := make(chan string)
	go markChannelDone(ch)

	for {
		time.Sleep(1 * time.Second)
		select {
		case v := <- ch:
			fmt.Println(v)
			return	
		default:
			fmt.Println("No value received")
		}
	}
}