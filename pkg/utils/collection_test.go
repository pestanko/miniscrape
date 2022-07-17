package utils

import (
	asrt "github.com/stretchr/testify/assert"
	"testing"
)

func TestIsInSlice(t *testing.T) {
	assert := asrt.New(t)

	needle := "needle"
	basicPred := func(s string) bool {
		return s == needle
	}
	tests := []struct {
		name     string
		haystack []string
		pred     func(string) bool
		result   bool
	}{
		{
			name:     "Empty haystack",
			haystack: []string{},
			pred:     basicPred,
			result:   false,
		},
		{
			name:     "Single needle in haystack",
			haystack: []string{needle},
			pred:     basicPred,
			result:   true,
		},
		{
			name:     "Multiple needles in haystack",
			haystack: []string{"foo", "bar", needle, "baz", needle},
			pred:     basicPred,
			result:   true,
		},
		{
			name:     "No needle in haystack",
			haystack: []string{"foo", "bar", "baz"},
			pred:     basicPred,
			result:   false,
		},
		{
			name:     "Needle at the end of haystack",
			haystack: []string{"foo", "bar", "baz", needle},
			pred:     basicPred,
			result:   true,
		},
		{
			name:     "Needle at the beginning of haystack",
			haystack: []string{needle, "foo", "bar", "baz"},
			pred:     basicPred,
			result:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert.Equal(test.result, IsInSlice(test.haystack, test.pred))
		})
	}
}
