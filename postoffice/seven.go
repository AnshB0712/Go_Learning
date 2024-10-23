package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Customer struct {
	Name string
	Item string
}

type DeskWorker struct {
	Name               string
	BackOfficeDeskChan chan string
}

func (w *DeskWorker) StartShift(queueChan chan Customer, shutdownChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case customer, ok := <-queueChan:
			if !ok {
				fmt.Printf("Desk worker %s is going home - queue closed\n", w.Name)
				return
			}
			time.Sleep(1 * time.Second)
			select {
			case w.BackOfficeDeskChan <- customer.Item:
				fmt.Printf("Desk worker %s sent %s to back office\n", w.Name, customer.Item)
			case <-shutdownChan:
				fmt.Printf("Desk worker %s is going home while processing\n", w.Name)
				return
			}
		case <-shutdownChan:
			fmt.Printf("Desk worker %s is going home\n", w.Name)
			return
		}
	}
}

type WizardBackOffice struct {
	Name string
}

func (w *WizardBackOffice) startDayShift(backOfficeDeskChan chan string, shutdownChan chan struct{}, processingWg *sync.WaitGroup) {
	defer processingWg.Done()
	
	// Create a separate WaitGroup for processing items
	var processWg sync.WaitGroup
	
	for {
		select {
		case item, ok := <-backOfficeDeskChan:
			if !ok {
				fmt.Printf("The back office channel is closed - waiting for remaining items, %s\n", w.Name)
				processWg.Wait() // Wait for all processing to complete
				return
			}
			fmt.Printf("Wizard: %s received %s\n", w.Name, item)
			processWg.Add(1)
			go func(item string) {
				defer processWg.Done()
				w.Process(item)
			}(item)
		case <-shutdownChan:
			fmt.Printf("The back office is closed - time to go home, %s\n", w.Name)
			processWg.Wait() // Wait for all processing to complete
			close(backOfficeDeskChan)
			return
		}
	}
}

func (w *WizardBackOffice) Process(item string) {
	fmt.Printf("Wizard %s's spell started processing %s...\n", w.Name, item)
	time.Sleep(10 * time.Second)
	fmt.Printf("Wizard %s's spell finished processing %s\n", w.Name, item)
}

func main() {
	start := time.Now()

	var (
		names = []string{
			"Alex",
			"Mia",
			"Juan",
			"Aisha",
			"Mohammad",
			"Isabella",
			"Ahmed",
			"Mei-Ling",
			"Leonardo",
			"Amara",
			"Rajesh",
			"Fatima",
			"Mateo",
			"Priya",
			"Carlos",
			"Lila",
			"Felix",
			"Gabriela",
			"Arjun",
			"Anika",
			"Giovanni",
			"Leila",
			"Manuel",
			"Isla",
			"Ali",
			"Lina",
			"Hugo",
			"Freya",
			"Javier",
			"Aylin",
			"Diego",
			"Emilia",
			"Ibrahim",
			"Yuki",
			"Aiden",
			"Elina",
			"Zhihao",
			"Anaya",
			"Mustafa",
			"Sienna",
			"Lily",
			"Amelie",
			"Maya",
			"Eva",
			"Oliver",
			"Samuel",
			"Liam",
			"Daniel",
			"Elijah",
			"Anna",
		}

		items = []string{
			"football",
			"box",
			"watermelon",
			"teddy bear",
			"basketball",
			"book",
			"gourmet chocolates",
			"holiday-themed socks",
			"personalized ornament",
			"miniature christmas tree",
			"holiday scented candles",
			"christmas-themed mug",
			"handmade soap set",
			"puzzle",
			"coffee sampler",
			"cozy knit scarf",
			"mini bottle of champagne",
			"essential oil diffuser",
			"festive cookie cutters",
			"mini photo album",
			"handwritten holiday card",
			"custom-made keychain",
			"a small plant",
			"pocket-sized board games",
			"holiday-themed puzzle",
			"popcorn seasoning kit",
			"mini holiday wreath",
			"wine sampler",
			"mini art supplies kit",
			"snow globe",
			"mini gingerbread house kit",
			"pocket-sized sketchbook",
			"festive face mask",
			"pocket-sized umbrella",
			"mini cheese and charcuterie board",
			"festive cocktail mixers",
			"mini holiday music cd",
			"handmade jewelry",
			"mini fairy lights",
			"miniature snowman kit",
			"funny holiday socks",
			"mini magnetic dartboard",
			"scratch-off world map",
			"mini photo frame",
			"reusable shopping bag",
			"mini hot sauce sampler",
			"pocket-sized tool kit",
			"mini bonsai tree kit",
			"holiday-themed coasters",
			"custom engraved keyring",
		}
	)

	queueChan := make(chan Customer)
	backOfficeDeskChan := make(chan string)
	shutdownChan := make(chan struct{})

	generateRandomCustomer := func() Customer {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		return Customer{
			Name: names[rand.Intn(len(names))],
			Item: items[rand.Intn(len(items))],
		}
	}

	var wg sync.WaitGroup

	w := WizardBackOffice{Name: "Tom Riddle"}
	dw := DeskWorker{Name: "Joe", BackOfficeDeskChan: backOfficeDeskChan}

	// Start the wizard
	wg.Add(1)
	go w.startDayShift(backOfficeDeskChan, shutdownChan, &wg)

	// Start the desk worker
	wg.Add(1)
	go dw.StartShift(queueChan, shutdownChan, &wg)

	// Start customer generator
	go func() {
		for {
			select {
			case <-shutdownChan:
				fmt.Println("Post office is closed - no more customers allowed")
				close(queueChan)
				return
			default:
				customer := generateRandomCustomer()
				fmt.Printf("%s enters the post office with %s\n", customer.Name, customer.Item)
				select {
				case queueChan <- customer:
				case <-shutdownChan:
					fmt.Println("Post office closed while customer was waiting")
					close(queueChan)
					return
				}
			}
		}
	}()

	// Trigger shutdown after 30 seconds
	time.AfterFunc(30*time.Second, func() {
		fmt.Println("30 seconds have passed - closing the post office")
		close(shutdownChan)
	})

	wg.Wait()
	fmt.Println("All workers have gone home ---------> ")

	end := time.Now()

	fmt.Printf("Time taken: %v\n", end.Sub(start))
}