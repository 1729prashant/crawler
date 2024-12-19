package main

import (
	"fmt"
	"net/url"
)

// crawlPage recursively crawls pages on the same domain as rawBaseURL
func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// Parse URLs to check if they're on the same domain
	base, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Error parsing base URL: %v\n", err)
		return
	}
	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error parsing current URL: %v\n", err)
		return
	}

	// Check if the current URL is on the same domain
	if base.Host != current.Host {
		fmt.Printf("Skipping external link: %s\n", rawCurrentURL)
		return
	}

	// Normalize URL
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing URL %s: %v\n", rawCurrentURL, err)
		return
	}

	// If we've seen this URL before, increment its count
	if count, exists := pages[normalizedURL]; exists {
		pages[normalizedURL] = count + 1
		fmt.Printf("Already visited %s, incrementing count to %d\n", normalizedURL, count+1)
		return
	}

	// Add new URL to pages map
	pages[normalizedURL] = 1
	fmt.Printf("New page found: %s\n", normalizedURL)

	// Get HTML content
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error fetching HTML for %s: %v\n", rawCurrentURL, err)
		return
	}

	// Get all URLs from this page
	urls, err := getURLsFromHTML(html, rawBaseURL)
	if err != nil {
		fmt.Printf("Error parsing URLs from %s: %v\n", rawCurrentURL, err)
		return
	}

	// Recursively crawl each link
	for _, link := range urls {
		crawlPage(rawBaseURL, link, pages)
	}
}
