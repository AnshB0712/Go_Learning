package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// PriceSummary represents the price information for a commodity
type PriceSummary struct {
	LastUpdated   string `json:"lastUpdated"`
	AvgPrice      string `json:"avgPrice"`
	LowestPrice   string `json:"lowestPrice"`
	HighestPrice  string `json:"highestPrice"`
}

// Option represents a select option from the HTML
type Option struct {
	Value string
	Text  string
}

// PriceMap stores the price data for all cities and commodities
type PriceMap struct {
	sync.RWMutex
	Data map[string]map[string]PriceSummary
}

// Progress tracks the scraping progress
type Progress struct {
	totalURLs     int32
	processed     int32
	successful    int32
	failed        int32
	startTime     time.Time
	logFrequency  time.Duration
	lastLogTime   time.Time
	mu            sync.Mutex
}

// HTTPClient with timeout
var client = &http.Client{
	Timeout: 5 * time.Second,
}

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

func main() {
	// Configure logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Println("Starting web scraper...")

	targetURL := "https://www.commodityonline.com/mandiprices/tomato/gujarat/surat"
	
	// Initialize price map with concurrent safe access
	priceMap := &PriceMap{
		Data: make(map[string]map[string]PriceSummary),
	}

	// Get initial options
	log.Println("Fetching initial options...")
	commodityOptions, marketOptions, err := extractOptionsFromURL(targetURL)
	if err != nil {
		log.Fatalf("Failed to extract options: %v", err)
	}

	log.Printf("Successfully fetched %d commodities and %d markets\n", len(commodityOptions), len(marketOptions))

	// Create URLs to visit
	urlsToVisit := generateURLs(commodityOptions, marketOptions)
	
	// Calculate total URLs
	totalURLs := 0
	for _, urls := range urlsToVisit {
		totalURLs += len(urls)
	}
	
	// Initialize progress tracker
	progress := &Progress{
		totalURLs:    int32(totalURLs),
		startTime:    time.Now(),
		logFrequency: 5 * time.Second,
		lastLogTime:  time.Now(),
	}

	log.Printf("Starting to process %d URLs...\n", totalURLs)

	// Process URLs concurrently
	processURLsConcurrently(urlsToVisit, priceMap, progress)

	// Final statistics
	duration := time.Since(progress.startTime)
	log.Printf("\nScraping completed in %v", duration.Round(time.Second))
	log.Printf("Total URLs processed: %d", atomic.LoadInt32(&progress.processed))
	log.Printf("Successful requests: %d", atomic.LoadInt32(&progress.successful))
	log.Printf("Failed requests: %d", atomic.LoadInt32(&progress.failed))

	// Save results
	savePriceMap(priceMap.Data)
}

func (p *Progress) logProgress() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if now.Sub(p.lastLogTime) >= p.logFrequency {
		processed := atomic.LoadInt32(&p.processed)
		successful := atomic.LoadInt32(&p.successful)
		failed := atomic.LoadInt32(&p.failed)
		percentage := float64(processed) / float64(p.totalURLs) * 100
		elapsed := time.Since(p.startTime)
		
		log.Printf("Progress: %.1f%% (%d/%d URLs) | Successful: %d | Failed: %d | Elapsed: %v",
			percentage, processed, p.totalURLs, successful, failed, elapsed.Round(time.Second))
		
		p.lastLogTime = now
	}
}

func processURLsConcurrently(urlsToVisit map[string][]string, priceMap *PriceMap, progress *Progress) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Limit concurrent requests

	for commodity, urls := range urlsToVisit {
		for _, url := range urls {
			wg.Add(1)
			go func(commodity, url string) {
				defer wg.Done()
				semaphore <- struct{}{} // Acquire semaphore
				defer func() { <-semaphore }() // Release semaphore

				exists, city := doesCommodityExistInCity(url)
				atomic.AddInt32(&progress.processed, 1)
				progress.logProgress()

				if exists {
					data := extractPriceSummary(url)
					if data.LastUpdated != "N/A" {
						atomic.AddInt32(&progress.successful, 1)
						log.Printf("Successfully scraped %s in %s\n", commodity, city)
					} else {
						atomic.AddInt32(&progress.failed, 1)
						log.Printf("Failed to scrape %s in %s\n", commodity, city)
					}
					
					priceMap.Lock()
					if _, ok := priceMap.Data[city]; !ok {
						priceMap.Data[city] = make(map[string]PriceSummary)
					}
					priceMap.Data[city][commodity] = data
					priceMap.Unlock()
				} else {
					atomic.AddInt32(&progress.failed, 1)
					log.Printf("Commodity %s not available in %s\n", commodity, city)
				}
			}(commodity, url)
		}
	}

	wg.Wait()
}

func extractPriceSummary(url string) PriceSummary {
	doc, err := fetchDocument(url)
	if err != nil {
		log.Printf("Error fetching document from %s: %v\n", url, err)
		return PriceSummary{
			LastUpdated:  "N/A",
			AvgPrice:    "N/A",
			LowestPrice: "N/A",
			HighestPrice: "N/A",
		}
	}

	summary := PriceSummary{}
	
	doc.Find(".mandi_highlight").Each(func(_ int, s *goquery.Selection) {
		summary.LastUpdated = strings.TrimSpace(strings.Replace(
			s.Find("p").First().Text(),
			"Last price updated: ",
			"",
			-1,
		))

		s.Find("h4").Each(func(_ int, h *goquery.Selection) {
			text := h.Text()
			value := strings.TrimSpace(h.Next().Text())
			
			switch text {
			case "Average Price":
				summary.AvgPrice = value
			case "Lowest Market Price":
				summary.LowestPrice = value
			case "Costliest Market Price":
				summary.HighestPrice = value
			}
		})
	})

	return summary
}

func doesCommodityExistInCity(url string) (bool, string) {
	doc, err := fetchDocument(url)
	if err != nil {
		log.Printf("Error checking commodity existence at %s: %v\n", url, err)
		return false, ""
	}

	// Check if commodity exists by looking for options
	exists := doc.Find("select[name='commodity'] option").Length() > 0
	city := extractCityFromURL(url)
	
	return exists, city
}

func extractOptionsFromURL(url string) ([]Option, []Option, error) {
	doc, err := fetchDocument(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch document: %v", err)
	}

	commodityOptions := extractSelectOptions(doc, "select[name='commodity'] option")
	marketOptions := extractSelectOptions(doc, "select[name='market'] option")

	return commodityOptions[2:], marketOptions[1:], nil
}

func extractSelectOptions(doc *goquery.Document, selector string) []Option {
	var options []Option
	
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		options = append(options, Option{
			Value: value,
			Text:  strings.TrimSpace(s.Text()),
		})
	})

	return options
}

func generateURLs(commodityOptions []Option, marketOptions []Option) map[string][]string {
	urls := make(map[string][]string)
	log.Printf("Generating URLs for %d commodities across %d markets...\n", 
		len(commodityOptions), len(marketOptions))
	
	for _, commodity := range commodityOptions {
		urls[commodity.Text] = make([]string, 0, len(marketOptions))
		for _, market := range marketOptions {
			url := fmt.Sprintf(
				"https://www.commodityonline.com/mandiprices/%s/gujarat/%s",
				commodity.Value,
				market.Value,
			)
			urls[commodity.Text] = append(urls[commodity.Text], url)
		}
	}

	return urls
}

func fetchDocument(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("User-Agent", userAgent)
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func extractCityFromURL(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func savePriceMap(data map[string]map[string]PriceSummary) {
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("price_map_%s.json", timestamp)
	
	file, err := os.Create(filepath.Join(".", filename))
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Fatalf("Failed to encode JSON: %v", err)
	}

	log.Printf("Price map saved to %s\n", filename)
}