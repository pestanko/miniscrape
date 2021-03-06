package filters

import (
	"strings"

	"github.com/pestanko/miniscrape/pkg/config"
	"github.com/rs/zerolog/log"
	"jaytaylor.com/html2text"
)

// NewHTMLConverter a new instance of the filter that
// uses the html2text converter
func NewHTMLConverter(page *config.Page) PageFilter {
	return &htmlFilterTags{
		page.Filters.HTML,
	}
}

type htmlFilterTags struct {
	html config.HTMLFilter
}

// Filter implements PageFilter
func (f *htmlFilterTags) Filter(content string) (string, error) {
	if f.html.Tables != "pretty" {
		content = useCustomHTMLTablesConverter(content)
	}

	text, err := html2text.FromString(content, html2text.Options{
		PrettyTables: f.html.Tables == "pretty",
		TextOnly:     f.html.TextOnly,
	})

	if err != nil {
		log.Warn().Err(err).Msg("Text extraction failed")
		return "", err
	}

	return text, nil
}

// IsEnabled implements PageFilter
func (*htmlFilterTags) IsEnabled() bool {
	return true
}

func (*htmlFilterTags) Name() string {
	return "html2text"
}

func useCustomHTMLTablesConverter(content string) string {
	if content == "" {
		return ""
	}

	content = strings.ReplaceAll(content, "<table", "<p")
	content = strings.ReplaceAll(content, "<TABLE", "<p")

	content = strings.ReplaceAll(content, "<tr", "<p")
	content = strings.ReplaceAll(content, "<TR", "<p")

	content = strings.ReplaceAll(content, "</tr>", "</p>")

	return strings.ReplaceAll(content, "</TR>", "</p>")
}
