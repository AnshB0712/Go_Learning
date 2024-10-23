package main

import (
	"fmt"
	"time"
	"sync"
	"math/rand"
)

// func calcSquares(n int, channel chan int) {
// 	sum := 0;
// 	for n != 0 {
// 		num := n % 10
// 		sum += (num * num)
// 		n = n / 10
// 	}
// 	channel <- sum
// }

// func calcCubes(n int, channel chan int) {
// 	sum := 0;
// 	for n != 0 {
// 		num := n % 10
// 		sum += (num * num * num)
// 		n = n / 10
// 	}
// 	// time.Sleep(4 * time.Second)
// 	channel <- sum
// }

// Write only channel for this function
func producer(ch chan <- int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}

	close(ch)
}

func writeToCh(ch chan <- int) {
	for i := 0; i < 5; i++ {
		ch <- i
		fmt.Println("Wrote: ", i)
	}

	close(ch)
}

// Worker Pool Implementation

type Job struct {
	id int
	random int
}

type Result struct {
	job Job
	sum int
}

var jobs = make(chan Job, 10)
var results = make(chan Result, 10)

func digitSums (n int) int {
	sum := 0
	num := n
	for  num != 0 {
		sum += num % 10
		num = num / 10
	}
	// time.Sleep(2 * time.Second)
	return sum
}

func worker(wg *sync.WaitGroup){
	defer wg.Done()
	for job := range jobs {
		output := Result{job, digitSums(job.random)}
		results <- output
	}
}

func createWorkerPool(workerCount int){
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	wg.Wait()
	close(results)
}

func allocate(noOfJobs int) {
	for i := 0; i < noOfJobs; i++ {
		random := rand.Intn(999)
		job := Job{id: i, random: random }
		jobs <- job
	}
	close(jobs)
}

func result(done chan bool) {
	for result := range results {
		fmt.Printf("Job id %d, input random no %d , sum of digits %d\n", result.job.id, result.job.random, result.sum)
	}
	done <- true
}

func main() {
	// sqChannel := make(chan int)
	// cuChannel := make(chan int)

	// go calcSquares(123, sqChannel)
	// go calcCubes(123, cuChannel)

	// squares := <-sqChannel
	// fmt.Println("Squares: ", squares)
	// cubes := <-cuChannel
	// fmt.Println("Cubes: ", cubes)

	// fmt.Println("Final output: ", squares + cubes)

	chanForInt := make(chan int)

	go producer(chanForInt)

	// for {
	// 	v, ok := <-chanForInt
	// 	if ok == false {
	// 		break
	// 	}
	// 	fmt.Println(v, ok)
	// }

	// better version for above loop
	for v := range chanForInt {
		fmt.Println(v)
	}

	// buffered channel read is blocking if its empty and write is blocking if its full
	// unbuffered channel, sending blocks until a receiver is ready to receive, and receiving blocks until a sender is ready to send.
	bufferedChannel := make(chan rune, 2)

	bufferedChannel <- 'A'
	bufferedChannel <- 'B'
	// bufferedChannel <- 'C' // this will block the program

	fmt.Println(string(<-bufferedChannel))
	fmt.Println(string(<-bufferedChannel))

	ch := make(chan int, 2)

	go writeToCh(ch)
	time.Sleep(2 * time.Second)
	for v := range ch {
		fmt.Println("Read: ", v)
		time.Sleep(2 * time.Second)
	}

	// Wait group for syncing go routines
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func (i int, wg *sync.WaitGroup) {
			fmt.Println("started Goroutine ", i)
			time.Sleep(2 * time.Second)
			fmt.Printf("Goroutine %d ended\n", i)
			wg.Done()
		}(i, &wg)
	}

	wg.Wait()



	// Worker Pool Implementation
	startTime := time.Now()
	noOfJobs := 100
	go allocate(noOfJobs)
	done := make(chan bool)
	go result(done)
	noOfWorkers := 10
	createWorkerPool(noOfWorkers)
	<-done
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")

}