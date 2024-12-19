package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// getHTML fetches the HTML content from the given URL.
func getHTML(rawURL string) (string, error) {
	// Fetch the URL
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for HTTP error status codes
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// Check if the content type is text/html
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("unexpected content type: %s", contentType)
	}

	// Read the body of the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
