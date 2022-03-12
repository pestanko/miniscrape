package scraper

import (
	"github.com/pestanko/miniscrape/scraper/config"
	"strings"
	"time"
)

type PageFilter interface {
	Filter(content string) (string, error)
	Enabled() bool
}

func NewCutFilter(page *config.Page) PageFilter {
	return &cutFilter{
		page.Filters.Cut,
	}
}
func NewDayFilter(page *config.Page) PageFilter {
	return &dayFilter{
		page.Filters.Day,
	}
}

type dayFilter struct {
	day config.DayFilter
}

func (f *dayFilter) Enabled() bool {
	return f.config().Enabled
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

func (f *cutFilter) Enabled() bool {
	return f.config().After != "" && f.config().Before != ""
}

func (f *cutFilter) Filter(content string) (string, error) {
	cfg := f.config()
	startIndex, endIndex := findBoundaries(content, cfg.Before, cfg.After)
	return cutContent(content, startIndex, endIndex), nil
}

func findBoundaries(content string, start string, end string) (int, int) {
	startIndex := -1
	if start != "" {
		startIndex = strings.LastIndex(content, start)
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

	return content[startIndex:endIndex]
}

func tryApplyDayFilter(content string, days []string, weekday time.Weekday) (int, int) {
	currIdx := (int(weekday) - 1) % 7
	nextIdx := (currIdx + 1) % 7
	var upperDays []string
	for _, day := range days {
		upperDays = append(upperDays, strings.ToUpper(day))
	}
	currDay := upperDays[currIdx]
	nextDay := upperDays[nextIdx]

	return findBoundaries(content, currDay, nextDay)
}
