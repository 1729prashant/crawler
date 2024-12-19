package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(urlString string) (string, error) {

	if urlString == "" {
		return "", fmt.Errorf("ERROR: empty string")
	}
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("ERROR: %s", err)
	}

	sanitisedURL := parsedURL.Host + strings.TrimRight(parsedURL.Path, "/")
	return sanitisedURL, nil
}
