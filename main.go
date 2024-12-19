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
}
