/*
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
		cfg.wg.Done()            // Ensure we decrement the WaitGroup even if an error or panic occurs
		<-cfg.concurrencyControl // Release the concurrency slot
	}()

	// Check if we've reached the maximum number of pages before any further processing
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		fmt.Println("Reached the maximum number of pages to crawl.")
		return // No need for cfg.wg.Done() here since it's in the defer
	}
	cfg.mu.Unlock()

	// Log when parsing the current URL
	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error parsing current URL %s: %v\n", rawCurrentURL, err)
		return // No need for cfg.wg.Done() here since it's in the defer
	}

	// Check if the current URL is on the same domain
	if cfg.baseURL.Host != current.Host {
		fmt.Printf("Skipping external link: %s\n", rawCurrentURL)
		return // No need for cfg.wg.Done() here since it's in the defer
	}

	// Normalize URL
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing URL %s: %v\n", rawCurrentURL, err)
		return // No need for cfg.wg.Done() here since it's in the defer
	}

	// Check if we've seen this URL before
	if !cfg.addPageVisit(normalizedURL) {
		fmt.Printf("Already visited %s, incrementing count\n", normalizedURL)
		return // No need for cfg.wg.Done() here since it's in the defer
	}

	fmt.Printf("New page found: %s\n", normalizedURL)

	// Get HTML content
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error fetching HTML for %s: %v\n", rawCurrentURL, err)
		return // No need for cfg.wg.Done() here since it's in the defer
	}

	// Get all URLs from this page
	urls, err := getURLsFromHTML(html, cfg.baseURL)
	if err != nil {
		fmt.Printf("Error parsing URLs from %s: %v\n", rawCurrentURL, err)
		return // No need for cfg.wg.Done() here since it's in the defer
	}

	// Check again before spawning new goroutines
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		fmt.Println("Reached the maximum number of pages to crawl during goroutine spawning.")
		return // No need for cfg.wg.Done() here since it's in the defer
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

		// Only spawn new goroutines if we haven't reached maxPages
		cfg.mu.Lock()
		if len(cfg.pages) >= cfg.maxPages {
			cfg.mu.Unlock()
			fmt.Println("Skipping further crawling due to maxPages limit.")
			continue
		}
		cfg.mu.Unlock()

		cfg.wg.Add(1) // Increment WaitGroup before spawning goroutine
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
*/

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
		cfg.wg.Done()            // Ensure we decrement the WaitGroup for this page
		<-cfg.concurrencyControl // Release the concurrency slot
	}()

	// Check if we've reached the maximum number of pages before any further processing
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		fmt.Println("Reached the maximum number of pages to crawl.")
		return
	}
	cfg.mu.Unlock()

	// Parse the current URL
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

	// Process found URLs
	for _, link := range urls {
		linkURL, err := url.Parse(link)
		if err != nil || linkURL.Host != cfg.baseURL.Host {
			fmt.Printf("Skipping invalid or external link: %s\n", link)
			continue
		}

		cfg.mu.Lock()
		if len(cfg.pages) >= cfg.maxPages {
			cfg.mu.Unlock()
			fmt.Println("Skipping further crawling due to maxPages limit.")
			return
		}

		// Try to acquire a concurrency slot without blocking
		select {
		case cfg.concurrencyControl <- struct{}{}:
			// Successfully acquired a slot, proceed with new goroutine
			cfg.wg.Add(1)
			go func(link string) {
				fmt.Printf("Goroutine for %s started\n", link)
				cfg.crawlPage(link)
			}(link)
		default:
			// If we can't acquire a slot, process synchronously
			fmt.Printf("Max concurrency reached, processing %s synchronously\n", link)
			cfg.mu.Unlock()
			continue
		}
		cfg.mu.Unlock()
	}
}
