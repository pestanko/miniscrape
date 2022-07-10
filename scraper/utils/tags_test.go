package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagsResolverForEmpty(t *testing.T) {
	assert := assert.New(t)

	resolver := MakeTagsResolver([]string{})

	assert.False(resolver.IsMatch([]string{"any"}))
}
