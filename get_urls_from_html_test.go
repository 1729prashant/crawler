package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
    <body>
        <a href="/path/one">
            <span>Boot.dev</span>
        </a>
        <a href="https://other.com/path/one">
            <span>Boot.dev</span>
        </a>
    </body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "only relative URLs",
			inputURL: "https://example.com",
			inputBody: `
<html>
    <body>
        <a href="/about">About</a>
        <a href="/contact">Contact</a>
    </body>
</html>
`,
			expected: []string{"https://example.com/about", "https://example.com/contact"},
		},
		{
			name:     "only absolute URLs",
			inputURL: "http://sample.com",
			inputBody: `
<html>
    <body>
        <a href="http://another.com/path">Another</a>
        <a href="https://third.com">Third</a>
    </body>
</html>
`,
			expected: []string{"http://another.com/path", "https://third.com"},
		},
		{
			name:     "nested links",
			inputURL: "https://nested.com",
			inputBody: `
<html>
    <body>
        <div>
            <a href="/home">Home</a>
            <p>
                <a href="https://external.com">External</a>
            </p>
        </div>
    </body>
</html>
`,
			expected: []string{"https://nested.com/home", "https://external.com"},
		},
		{
			name:     "mixed URLs with no <a> tags",
			inputURL: "https://noanchor.com",
			inputBody: `
<html>
    <body>
        <img src="/images/logo.png">
        <script src="https://cdn.example.com/script.js"></script>
    </body>
</html>
`,
			expected: []string{"https://noanchor.com/images/logo.png", "https://cdn.example.com/script.js"},
		},
		{
			name:      "empty input",
			inputURL:  "https://empty.com",
			inputBody: ``,
			expected:  []string{},
		},
		{
			name:      "malformed HTML",
			inputURL:  "https://malformed.com",
			inputBody: `<a href="malformed">Link</a`,
			expected:  []string{"https://malformed.com/malformed"}, // This might work if parser is lenient
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the input URL to match the function signature
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Fatalf("Failed to parse base URL for test case %s: %v", tc.name, err)
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Errorf("%s: unexpected error: %v", tc.name, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("%s: expected URLs %v, but got %v", tc.name, tc.expected, actual)
			}
		})
	}
}
