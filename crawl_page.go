package main

import (
	"fmt"
	"net/url"
)

// crawlPage recursively crawls pages on the same domain as rawBaseURL
func (cfg *config) crawlPage(rawCurrentURL string) {
	fmt.Printf("Starting crawl of: %s\n", rawCurrentURL)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
		fmt.Printf("Completed crawl of: %s\n", rawCurrentURL)
		cfg.wg.Done()
		<-cfg.concurrencyControl // Release the slot in the channel
	}()

	// Check if we've reached the maximum number of pages before any further processing
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		fmt.Println("Reached the maximum number of pages to crawl.")
		return
	}
	cfg.mu.Unlock()

	// Log when parsing the current URL
	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error parsing current URL %s: %v\n", rawCurrentURL, err)
		return
	}

	// Check if the current URL is on the same domain
	if cfg.baseURL.Host != current.Host {
		fmt.Printf("Skipping external link: %s\n", rawCurrentURL)
		return
	}

	// Normalize URL
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing URL %s: %v\n", rawCurrentURL, err)
		return
	}

	// Check if we've seen this URL before
	if !cfg.addPageVisit(normalizedURL) {
		fmt.Printf("Already visited %s, incrementing count\n", normalizedURL)
		return
	}

	fmt.Printf("New page found: %s\n", normalizedURL)

	// Get HTML content
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error fetching HTML for %s: %v\n", rawCurrentURL, err)
		return
	}

	// Get all URLs from this page
	urls, err := getURLsFromHTML(html, cfg.baseURL)
	if err != nil {
		fmt.Printf("Error parsing URLs from %s: %v\n", rawCurrentURL, err)
		return
	}

	// Check again before spawning new goroutines to ensure we haven't exceeded maxPages due to concurrent updates
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		fmt.Println("Reached the maximum number of pages to crawl during goroutine spawning.")
		return
	}
	cfg.mu.Unlock()

	// Log before spawning new goroutines
	fmt.Printf("Spawning %d new goroutines for %s\n", len(urls), rawCurrentURL)
	for _, link := range urls {
		linkURL, err := url.Parse(link)
		if err != nil || linkURL.Host != cfg.baseURL.Host {
			fmt.Printf("Skipping invalid or external link: %s\n", link)
			continue // Skip this link
		}

		cfg.wg.Add(1) // This stays outside - we must increment BEFORE spawning goroutine
		go func(link string) {
			fmt.Printf("Goroutine for %s started\n", link)
			defer func() {
				fmt.Printf("Goroutine for %s completed\n", link)
				if r := recover(); r != nil {
					fmt.Println("Recovered from panic in goroutine:", r)
				}
				<-cfg.concurrencyControl
				cfg.wg.Done()
			}()
			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(link)
		}(link)
	}
}
