package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

// Helper method to add a page visit, ensuring thread-safety and checking if it's the first visit
func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, exists := cfg.pages[normalizedURL]; !exists {
		cfg.pages[normalizedURL] = 1
		return true
	}
	cfg.pages[normalizedURL]++
	return false
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: ./crawler <baseURL> <maxPages> <maxConcurrency>")
		os.Exit(1)
	}

	baseURL := os.Args[1]
	maxPages, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid number for max pages:", os.Args[2])
		os.Exit(1)
	}
	maxConcurrency, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid number for max concurrency:", os.Args[3])
		os.Exit(1)
	}

	fmt.Printf("Starting crawl of: %s with max pages %d and max concurrency %d\n", baseURL, maxPages, maxConcurrency)

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Printf("Error parsing base URL: %v\n", err)
		os.Exit(1)
	}

	// Initialize config
	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	// Start crawling from the base URL
	cfg.wg.Add(1)
	cfg.concurrencyControl <- struct{}{}
	go cfg.crawlPage(baseURL)

	// Wait for all goroutines to finish
	cfg.wg.Wait()

	// Print the results
	printReport(cfg.pages, baseURL)

}
