package filters

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/pestanko/miniscrape/scraper/config"
)

// NewHTMLToMdConverter a new instance of the filter that
// converts html to markdown
func NewHTMLToMdConverter(page *config.Page) PageFilter {
	return &htmlToMdConverter{
		page.Filters.HTML,
	}
}

type htmlToMdConverter struct {
	html config.HTMLFilter
}

// Filter implements PageFilter
func (*htmlToMdConverter) Filter(content string) (string, error) {
	converter := makeMdConverter()

	return converter.ConvertString(content)
}

// IsEnabled implements PageFilter
func (*htmlToMdConverter) IsEnabled() bool {
	return true
}

// Name implements PageFilter
func (*htmlToMdConverter) Name() string {
	return "html2md"
}

func makeMdConverter() *md.Converter {
	converter := md.NewConverter("", true, nil)
	// Use the `GitHubFlavored` plugin from the `plugin` package.
	converter.Use(plugin.GitHubFlavored())

	return converter
}
