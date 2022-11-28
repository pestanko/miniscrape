package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagsResolverForEmpty(t *testing.T) {
	s := assert.New(t)

	resolver := MakeTagsResolver([]string{})

	s.True(resolver.IsMatch([]string{"any"}))
}
