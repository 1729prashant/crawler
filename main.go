package main

import (
	"fmt"
	"os"
)

func main() {
	// Check if exactly one command-line argument (besides the program name) is provided
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	// If we have exactly one argument, it's the BASE_URL
	baseURL := os.Args[1]
	fmt.Printf("starting crawl of: %s\n", baseURL)

	// Fetch HTML
	html, err := getHTML(baseURL)
	if err != nil {
		fmt.Printf("Error fetching HTML: %v\n", err)
		os.Exit(1)
	}

	// Print some HTML (first 500 characters for brevity)
	fmt.Println("HTML:")
	fmt.Println(html) // Print only first 5000 characters to avoid overwhelming output
}
