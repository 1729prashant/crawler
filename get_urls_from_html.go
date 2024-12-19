package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// getURLsFromHTML parses the given HTML body to find all URLs in <a>, <img>, and <script> tags,
// converting relative URLs to absolute URLs using the rawBaseURL.
func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	// Early return for empty HTML body to avoid parsing errors
	if htmlBody == "" {
		return []string{}, nil
	}

	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var links []string
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a", "img", "script":
				for _, attr := range n.Attr {
					var attrKey string
					switch n.Data {
					case "a":
						attrKey = "href"
					case "img", "script":
						attrKey = "src"
					}
					if attr.Key == attrKey {
						linkURL, err := url.Parse(attr.Val)
						if err != nil {
							continue // Skip this link if URL is invalid
						}
						absoluteURL := baseURL.ResolveReference(linkURL).String()
						links = append(links, absoluteURL)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}
	findLinks(doc)
	return links, nil
}
