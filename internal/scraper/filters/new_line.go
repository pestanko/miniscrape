package filters

import (
	"regexp"

	"github.com/pestanko/miniscrape/internal/models"
)

var normPattern = regexp.MustCompile("[\n]+")

// NewNewLineTrimConverter a new instance of the filter that
// cuts the line of the content
func NewNewLineTrimConverter(_ *models.Page) PageFilter {
	return &newLineTrimConverter{}
}

type newLineTrimConverter struct{}

// Filter implements PageFilter
func (*newLineTrimConverter) Filter(content string) (string, error) {
	return normPattern.ReplaceAllString(content, "\n"), nil
}

// IsEnabled implements PageFilter
func (*newLineTrimConverter) IsEnabled() bool {
	return true
}

// Name implements PageFilter
func (*newLineTrimConverter) Name() string {
	return "newline"
}
