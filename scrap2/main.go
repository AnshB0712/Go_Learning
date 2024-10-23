package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/time/rate"
)

// PriceMap stores the scraped data
type PriceMap struct {
	sync.RWMutex
	Data map[string]map[string]PriceSummary
}

// PriceSummary represents the price information for a commodity
type PriceSummary struct {
	LastUpdated  string `json:"lastUpdated"`
	AvgPrice     string `json:"avgPrice"`
	LowestPrice  string `json:"lowestPrice"`
	HighestPrice string `json:"highestPrice"`
}

// Option represents a select option from the HTML
type Option struct {
	Value string
	Text  string
}

// ScrapeJob represents a single URL to be scraped
type ScrapeJob struct {
	URL       string
	Commodity string
}

// Client wraps http.Client with custom configuration and rate limiting
type Client struct {
	http.Client
	limiter *rate.Limiter
}

func newClient(rps float64) *Client {
	return &Client{
		Client: http.Client{
			Timeout: 5 * time.Second,
		},
		limiter: rate.NewLimiter(rate.Limit(rps), 1), // Allow burst of 1
	}
}

func (c *Client) fetch(url string) (*goquery.Document, error) {
	// Wait for rate limiter
	err := c.limiter.Wait(context.Background())
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	return doc, nil
}

func (pm *PriceMap) updatePrice(city, commodity string, summary PriceSummary) {
	pm.Lock()
	defer pm.Unlock()

	if _, exists := pm.Data[city]; !exists {
		pm.Data[city] = make(map[string]PriceSummary)
	}
	pm.Data[city][commodity] = summary
}

// Function to extract city from URL
func extractCityFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 6 {
		return parts[5]
	}
	return ""
}

// Function to extract price summary from HTML document
func extractPriceSummary(doc *goquery.Document) PriceSummary {
	summary := PriceSummary{}

	doc.Find(".mandi_highlight").Each(func(_ int, s *goquery.Selection) {
		summary.LastUpdated = strings.TrimPrefix(
			s.Find("p").First().Text(),
			"Last price updated: ",
		)

		s.Find("h4").Each(func(_ int, h *goquery.Selection) {
			text := h.Text()
			switch text {
			case "Average Price":
				summary.AvgPrice = strings.TrimSpace(h.Next().Text())
			case "Lowest Market Price":
				summary.LowestPrice = strings.TrimSpace(h.Next().Text())
			case "Costliest Market Price":
				summary.HighestPrice = strings.TrimSpace(h.Next().Text())
			}
		})
	})

	fmt.Printf("Extracted summary: %+v\n", summary)

	return summary
}

// Function to extract options from the document
func extractOptions(doc *goquery.Document, selector string, skip int) []Option {
	var options []Option
	doc.Find(selector).Find("option").Each(func(i int, s *goquery.Selection) {
		if i < skip {
			return
		}
		options = append(options, Option{
			Value: s.AttrOr("value", ""),
			Text:  strings.TrimSpace(s.Text()),
		})
	})
	return options
}

// Function to save the price map to a JSON file
func savePriceMap(priceMap *PriceMap) error {
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("price_map_%s.json", timestamp)

	priceMap.RLock()
	data, err := json.MarshalIndent(priceMap.Data, "", "  ")
	priceMap.RUnlock()

	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	err = os.WriteFile(filepath.Join(".", filename), data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	fmt.Printf("Price map saved to %s\n", filename)
	return nil
}

// Main concurrent URL processing function
func processURLsConcurrently(client *Client, urlsToVisit map[string][]string) (*PriceMap, error) {
	priceMap := &PriceMap{
		Data: make(map[string]map[string]PriceSummary),
	}

	// Channel for errors
	errorsChan := make(chan error, len(urlsToVisit)*10)
	var wg sync.WaitGroup

	// Iterate over the URLs and spawn goroutines for each job
	for commodity, urls := range urlsToVisit {
		for _, url := range urls {
			wg.Add(1) // Add a job to the wait group
			go func(commodity, url string) {
				defer wg.Done()

				fmt.Printf("Processing %s for %s\n", url, commodity)

				// Retry mechanism inside the goroutine
				var doc *goquery.Document
				var err error
				for retries := 0; retries < 1; retries++ {
					doc, err = client.fetch(url)
					if err == nil {
						break
					}
					fmt.Printf("Retry %d for %s\n", retries+1, url)
					time.Sleep(1 * time.Second) // Delay before retry
				}

				if err != nil {
					fmt.Printf("Error fetching %s: %v\n", url, err)
					errorsChan <- fmt.Errorf("error fetching %s: %w", url, err)
					return
				}

				city := extractCityFromURL(url)
				if doc != nil {
					summary := extractPriceSummary(doc)
					priceMap.updatePrice(city, commodity, summary)
				} else {
					priceMap.updatePrice(city, commodity, PriceSummary{
						LastUpdated:  "N/A",
						AvgPrice:     "N/A",
						LowestPrice:  "N/A",
						HighestPrice: "N/A",
					})
				}
			}(commodity, url)
		}
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errorsChan)
	}()

	// Collect errors
	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		fmt.Printf("Encountered %d errors during scraping, continuing...\n", len(errors))
	}

	return priceMap, nil
}

func main() {
	// Configure client with rate limiting (5 requests per second)
	client := newClient(5.0)
	targetURL := "https://www.commodityonline.com/mandiprices/tomato/gujarat/surat"

	fmt.Println("Fetching options from target URL...")
	doc, err := client.fetch(targetURL)
	if err != nil {
		fmt.Printf("Error fetching target URL: %v\n", err)
		os.Exit(1)
	}

	commodityOptions := extractOptions(doc, "select[name=\"commodity\"]", 2)
	marketOptions := extractOptions(doc, "select[name=\"market\"]", 1)

	fmt.Printf("Fetched %d commodities and %d markets\n", len(commodityOptions), len(marketOptions))

	// Create commodity and market maps
	commodityMap := make(map[string]string)
	for _, opt := range commodityOptions {
		commodityMap[opt.Text] = opt.Value
	}

	marketMap := make(map[string]string)
	for _, opt := range marketOptions {
		marketMap[opt.Text] = opt.Value
	}

	// Generate URLs to visit
	urlsToVisit := make(map[string][]string)
	for _, commodity := range commodityOptions {
		urls := make([]string, 0, len(marketOptions))
		for _, market := range marketOptions {
			url := fmt.Sprintf(
				"https://www.commodityonline.com/mandiprices/%s/gujarat/%s",
				commodityMap[commodity.Text],
				marketMap[market.Text],
			)
			urls = append(urls, url)
		}
		urlsToVisit[commodity.Text] = urls
	}

	// Calculate optimal number of workers based on CPU cores
	numWorkers := runtime.NumCPU() * 5 // Use 2 workers per CPU core
	fmt.Printf("Starting to process URLs with %d workers...\n", numWorkers)

	startTime := time.Now()
	priceMap, err := processURLsConcurrently(client, urlsToVisit)
	if err != nil {
		fmt.Printf("Error processing URLs: %v\n", err)
		os.Exit(1)
	}

	if err := savePriceMap(priceMap); err != nil {
		fmt.Printf("Error saving price map: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(startTime)
	fmt.Printf("Script executed successfully in %v\n", duration)
}
