package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
		wantErr  bool
	}{
		{
			name:     "remove https scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
			wantErr:  false,
		},
		{
			name:     "remove https scheme, remove end slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
			wantErr:  false,
		},
		{
			name:     "remove http scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
			wantErr:  false,
		},
		{
			name:     "remove http scheme, remove end slash",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
			wantErr:  false,
		},
		// Additional test cases:
		{
			name:     "only domain with end slash",
			inputURL: "http://example.com/",
			expected: "example.com",
			wantErr:  false,
		},
		{
			name:     "only domain",
			inputURL: "http://example.com",
			expected: "example.com",
			wantErr:  false,
		},
		{
			name:     "query parameters",
			inputURL: "http://example.com/path?query=param",
			expected: "example.com/path",
			wantErr:  false,
		},
		{
			name:     "fragment identifier",
			inputURL: "http://example.com/path#fragment",
			expected: "example.com/path",
			wantErr:  false,
		},
		{
			name:     "invalid URL",
			inputURL: "://example.com",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "empty URL",
			inputURL: "",
			expected: "",
			wantErr:  true,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if (err != nil) != tc.wantErr {
				t.Errorf("Test %v - '%s' FAIL: error mismatch: wantErr %v, got %v", i, tc.name, tc.wantErr, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
