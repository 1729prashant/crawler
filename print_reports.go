package main

import (
	"fmt"
	"sort"
	"strings"
)

// pageCount represents a page and its link count
type pageCount struct {
	url   string
	count int
}

// sortPageCounts sorts the pages map into a slice of pageCount structs
// ordered by count (descending) and URL (ascending) for equal counts
func sortPageCounts(pages map[string]int) []pageCount {
	// Convert map to slice of pageCount structs
	pageCounts := make([]pageCount, 0, len(pages))
	for url, count := range pages {
		pageCounts = append(pageCounts, pageCount{url: url, count: count})
	}

	// Sort slice by count (descending) and URL (ascending) for equal counts
	sort.Slice(pageCounts, func(i, j int) bool {
		if pageCounts[i].count != pageCounts[j].count {
			return pageCounts[i].count > pageCounts[j].count // Higher counts first
		}
		return pageCounts[i].url < pageCounts[j].url // Alphabetical for equal counts
	})

	return pageCounts
}

// printReport prints a formatted report of the crawl results
func printReport(pages map[string]int, baseURL string) {
	// Print header
	header := fmt.Sprintf("REPORT for %s", baseURL)
	separator := strings.Repeat("=", len(header)+4) // +4 for padding

	fmt.Printf("\n%s\n  %s\n%s\n\n", separator, header, separator)

	// Sort and print pages
	sortedPages := sortPageCounts(pages)
	for _, pc := range sortedPages {
		fmt.Printf("Found %d internal links%s to %s\n",
			pc.count,
			pluralize(pc.count), // Add 's' for counts other than 1
			pc.url)
	}
	fmt.Println() // Add final newline for cleaner output
}

// pluralize returns "s" if count is not 1, empty string otherwise
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
