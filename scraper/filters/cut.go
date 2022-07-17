package filters

import (
	"strings"

	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/rs/zerolog/log"
)

// NewCutFilter create a new instance of the cut filter
func NewCutFilter(page *config.Page) PageFilter {
	return &cutFilter{
		page.Filters.Cut,
	}
}

// NewCutLineFilter create a new instance of the cut line filter
func NewCutLineFilter(page *config.Page) PageFilter {
	return &cutLineFilter{
		page.Filters.CutLine,
	}
}

type cutFilter struct {
	cut config.CutFilter
}

func (f *cutFilter) config() *config.CutFilter {
	return &f.cut
}

func (f *cutFilter) IsEnabled() bool {
	return f.config().After != "" || f.config().Before != ""
}

func (*cutFilter) Name() string {
	return "cut"
}

func (f *cutFilter) Filter(content string) (string, error) {
	cfg := f.config()
	startIndex, endIndex := findBoundaries(content, cfg.Before, cfg.After)
	return cutContent(content, startIndex, endIndex), nil
}

type cutLineFilter struct {
	cutLine config.CutLineFilter
}

func (c *cutLineFilter) Filter(content string) (string, error) {
	cfg := c.config()
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		line := line
		if cfg.Contains != "" && strings.Contains(line, cfg.Contains) {
			line = ""
		}
		if cfg.StartsWith != "" && strings.HasPrefix(line, cfg.StartsWith) {
			line = ""
		}
		if cfg.CutAfter != "" {
			start, end := findBoundaries(line, "", cfg.CutAfter)
			line = cutContent(line, start, end)
		}
		if line != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n"), nil
}

func (c *cutLineFilter) IsEnabled() bool {
	cfg := c.config()
	return cfg.Contains != "" || cfg.CutAfter != "" || cfg.StartsWith != ""
}

func (c *cutLineFilter) config() *config.CutLineFilter {
	return &c.cutLine
}

func (*cutLineFilter) Name() string {
	return "cut_line"
}

func cutContent(content string, startIndex int, endIndex int) string {
	if startIndex == -1 {
		startIndex = 0
	}

	if endIndex == -1 {
		endIndex = len(content) - 1
	}

	if content == "" {
		return ""
	}

	log.Debug().
		Int("from", startIndex).
		Int("to", endIndex).
		Msg("Content range")

	return content[startIndex:endIndex]
}

func findBoundaries(content string, start string, end string) (int, int) {
	startIndex := -1
	if start != "" {
		startIndex = strings.Index(content, start)
	}

	if startIndex == -1 {
		startIndex = 0
	}

	endIndex := -1
	if end != "" {
		endIndex = strings.Index(content, end)
	}

	return startIndex, endIndex
}
