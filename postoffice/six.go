package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type WizardBackOffice struct {
	Name string
}

func (w *WizardBackOffice) startDayShift(backOfficeDeskChan chan string, shutdownChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case item, ok := <-backOfficeDeskChan:
			if ok {
				fmt.Printf("WIZARD: %s received %s\n", w.Name, item)

				go func(item string) {
					fmt.Printf("Wizard %s casted a spell to process %s\n", w.Name, item)
					w.Process(item)
				}(item)
			}
		case <-shutdownChan:
			fmt.Printf("The back office is closed - time to go home, %s\n", w.Name)
			return
		}
	}
}

func (w *WizardBackOffice) Process(item string) {
	fmt.Printf("Wizard %s's spell started processing %s...\n", w.Name, item)

	// to simulate long processing
	time.Sleep(10 * time.Second)

	fmt.Printf("Wizard %s's spell finished processing %s\n", w.Name, item)
}

type Customer struct {
	Name string
	Item string
}

type DeskWorker struct {
	Name              string
	BackOfficeDeskChan chan string
}

func (w *DeskWorker) startDayShift(deskChan chan string, phoneChan chan string, shutdownChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case item, ok := <-deskChan:
			if ok {
				time.Sleep(1 * time.Second)
				w.BackOfficeDeskChan <- item
				fmt.Printf("Desk worker %s started passing %s to the back office...\n", w.Name, item)
				fmt.Printf("Desk worker %s passed %s to the back office\n\n", w.Name, item)
			}
		case call, ok := <-phoneChan:
			if ok {
				fmt.Printf("CALL --->  %s\n", call)
			}
		case <-shutdownChan:
			fmt.Printf("The desk is closed - time to go home, %s\n", w.Name)
			return
		}
	}
}

func (w *DeskWorker) Receive(item string) {
	fmt.Printf("Deskworker: %s has 'received' the parcel containing: %s\n", w.Name, item)
	fmt.Printf("Deskworker: %s has 'started checking ID of the customer' for the parcel containing: %s\n", w.Name, item)
}

func main() {

	start := time.Now()
	var (
		names = []string{
			"Alex", "Mia", "Juan", "Aisha", "Mohammad", "Isabella", "Ahmed", "Mei-Ling", "Leonardo", "Amara",
			"Rajesh", "Fatima", "Mateo", "Priya", "Carlos", "Lila", "Felix", "Gabriela", "Arjun", "Anika",
			"Giovanni", "Leila", "Manuel", "Isla", "Ali", "Lina", "Hugo", "Freya", "Javier", "Aylin",
			"Diego", "Emilia", "Ibrahim", "Yuki", "Aiden", "Elina", "Zhihao", "Anaya", "Mustafa", "Sienna",
			"Lily", "Amelie", "Maya", "Eva", "Oliver", "Samuel", "Liam", "Daniel", "Elijah", "Anna",
		}

		items = []string{
			"football", "box", "watermelon", "teddy bear", "basketball", "book", "gourmet chocolates", "holiday-themed socks",
			"personalized ornament", "miniature christmas tree", "holiday scented candles", "christmas-themed mug", "handmade soap set", "puzzle",
			"coffee sampler", "cozy knit scarf", "mini bottle of champagne", "essential oil diffuser", "festive cookie cutters", "mini photo album",
			"handwritten holiday card", "custom-made keychain", "a small plant", "pocket-sized board games", "holiday-themed puzzle", "popcorn seasoning kit",
			"mini holiday wreath", "wine sampler", "mini art supplies kit", "snow globe", "mini gingerbread house kit", "pocket-sized sketchbook",
			"festive face mask", "pocket-sized umbrella", "mini cheese and charcuterie board", "festive cocktail mixers", "mini holiday music cd",
			"handmade jewelry", "mini fairy lights", "miniature snowman kit", "funny holiday socks", "mini magnetic dartboard", "scratch-off world map",
			"mini photo frame", "reusable shopping bag", "mini hot sauce sampler", "pocket-sized tool kit", "mini bonsai tree kit", "holiday-themed coasters",
			"custom engraved keyring",
		}
	)

	generateRandomCustomer := func() Customer {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		return Customer{
			Name: names[rand.Intn(len(names))],
			Item: items[rand.Intn(len(items))],
		}
	}

	deskChan := make(chan string, 3)
	phoneChan := make(chan string, 1)
	queueChan := make(chan Customer)
	backOfficeDeskChan := make(chan string, 3)

	deskShutdownChan := make(chan struct{})
	postOfficeShutdownChan := make(chan struct{})

	wizard := WizardBackOffice{Name: "Dumbledore"}
	bobWorker := DeskWorker{Name: "Bob", BackOfficeDeskChan: backOfficeDeskChan}

	// Add wait groups to ensure proper shutdown
	var wg sync.WaitGroup
	wg.Add(2) // One for desk worker, one for wizard back office

	go func() {
		wizard.startDayShift(backOfficeDeskChan, postOfficeShutdownChan, &wg)
	}()

	go func() {
		bobWorker.startDayShift(deskChan, phoneChan, deskShutdownChan, &wg)
	}()

	// Customer queue generation
	go func() {
		for {
			select {
			case <-postOfficeShutdownChan:
				fmt.Println("The post office is closed - time to go home")
				close(queueChan)
				return
			default:
				customer := generateRandomCustomer()
				fmt.Printf("%s enters the post office with gift: %v\n", customer.Name, customer.Item)
				queueChan <- customer
			}
		}
	}()

	go func() {
		phoneChan <- "package arrived yet?"
		time.Sleep(8 * time.Second)
		phoneChan <- "now?"
	}()

	time.AfterFunc(5*time.Second, func() {
		close(postOfficeShutdownChan)
	})

	for c := range queueChan {
		bobWorker.Receive(c.Item)
		deskChan <- c.Item
		fmt.Printf("Customer: %s has 'queued' the parcel: %s\n", c.Name, c.Item)
	}

	// Close the desk and back office channels after shutdown
	close(deskShutdownChan)
	close(deskChan)
	wg.Wait() // Wait for all goroutines to finish

	end := time.Now()

	fmt.Println("------------------Time_Taken: --------------------", end.Sub(start).Seconds())
}
