package scraper

import (
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"

	"github.com/pestanko/miniscrape/scraper/config"

	"jaytaylor.com/html2text"
)

type PageFilter interface {
	Filter(content string) (string, error)
	IsEnabled() bool
	Name() string
}

func NewCutFilter(page *config.Page) PageFilter {
	return &cutFilter{
		page.Filters.Cut,
	}
}

func NewCutLineFilter(page *config.Page) PageFilter {
	return &cutLineFilter{
		page.Filters.CutLine,
	}
}

func NewDayFilter(page *config.Page) PageFilter {
	return &dayFilter{
		page.Filters.Day,
	}
}

func NewHTMLConverter(page *config.Page) PageFilter {
	return &htmlFilterTags{
		page.Filters.Html,
	}
}

func NewHTMLToMdConverter(page *config.Page) PageFilter {
	return &htmlToMdConverter{
		page.Filters.Html,
	}
}

func NewNewLineTrimConverter(page *config.Page) PageFilter {
	return &newLineTrimConverter{}
}

type dayFilter struct {
	day config.DayFilter
}

func (f *dayFilter) IsEnabled() bool {
	return f.config().Enabled
}

func (*dayFilter) Name() string {
	return "day"
}

func (f *dayFilter) config() *config.DayFilter {
	return &f.day
}

func (f *dayFilter) Filter(content string) (string, error) {
	days := f.config().Days
	weekday := time.Now().Weekday()
	upperContent := strings.ToUpper(content)
	if len(days) != 0 {
		start, end := tryApplyDayFilter(upperContent, days, weekday)
		return cutContent(content, start, end), nil
	} else {
		allVersions := [][]string{
			{"Pondělí", "Úterý", "Středa", "Čtvrtek", "Pátek", "Sobota", "Neděle"},
			{"Pondeli", "Uteri", "Streda", "Ctvrtek", "Patek", "Sobota", "Nedele"},
			{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		}
		for _, days := range allVersions {
			start, end := tryApplyDayFilter(upperContent, days, weekday)
			if start == -1 && end == -1 {
				continue
			}
			return cutContent(content, start, end), nil
		}
	}
	return content, nil
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

func tryApplyDayFilter(content string, days []string, weekday time.Weekday) (int, int) {
	currIdx := (int(weekday) - 1) % 7
	if currIdx < 0 {
		currIdx = 6
	}
	nextIdx := (currIdx + 1) % 7
	var upperDays []string
	for _, day := range days {
		upperDays = append(upperDays, strings.ToUpper(day))
	}
	currDay := upperDays[currIdx]
	nextDay := upperDays[nextIdx]

	return findBoundaries(content, currDay, nextDay)
}

type htmlFilterTags struct {
	html config.HtmlFilter
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

type htmlToMdConverter struct {
	html config.HtmlFilter
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
