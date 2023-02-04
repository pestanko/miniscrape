package filters

import (
	"github.com/pestanko/miniscrape/internal/config"
	"strings"
	"time"
)

// NewDayFilter a new instance of the filter that
// cuts a content based on days
func NewDayFilter(page *config.Page) PageFilter {
	return &dayFilter{
		page.Filters.Day,
	}
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
	}
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
	return content, nil
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
