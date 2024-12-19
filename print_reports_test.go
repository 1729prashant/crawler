package main

import (
	"reflect"
	"testing"
)

func TestSortPageCounts(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []pageCount
	}{
		{
			name: "Sort by count descending",
			input: map[string]int{
				"page1": 1,
				"page2": 3,
				"page3": 2,
			},
			expected: []pageCount{
				{url: "page2", count: 3},
				{url: "page3", count: 2},
				{url: "page1", count: 1},
			},
		},
		{
			name: "Sort alphabetically for equal counts",
			input: map[string]int{
				"pageB": 2,
				"pageA": 2,
				"pageC": 2,
			},
			expected: []pageCount{
				{url: "pageA", count: 2},
				{url: "pageB", count: 2},
				{url: "pageC", count: 2},
			},
		},
		{
			name: "Mixed sorting",
			input: map[string]int{
				"pageB": 2,
				"pageA": 3,
				"pageD": 2,
				"pageC": 1,
			},
			expected: []pageCount{
				{url: "pageA", count: 3},
				{url: "pageB", count: 2},
				{url: "pageD", count: 2},
				{url: "pageC", count: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortPageCounts(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("sortPageCounts() = %v, want %v", result, tt.expected)
			}
		})
	}
}
