package web

import (
	"net/http"

	"github.com/pestanko/miniscrape/scraper"
	"github.com/pestanko/miniscrape/scraper/config"
)

// HandleCategories handler
func HandleCategories(
	scrapeService *scraper.Service,
	w http.ResponseWriter,
	_ *http.Request,
) {
	var dto []categoryDto
	for _, cat := range scrapeService.GetCategories() {
		catDto := categoryDto{
			Name:  cat.Name,
			Tags:  getAllTagsForCategory(cat.Pages),
			Pages: getPagesForCategory(cat.Pages),
		}
		dto = append(dto, catDto)
	}

	WriteJSONResponse(w, http.StatusOK, dto)
}

func getPagesForCategory(pages []config.Page) []string {
	var pageNames = make([]string, len(pages))

	for i, pg := range pages {
		pageNames[i] = pg.CodeName
	}

	return pageNames
}

func getAllTagsForCategory(pages []config.Page) []string {
	tagsSet := make(map[string]bool)
	for _, pg := range pages {
		for _, pageTag := range pg.Tags {
			tagsSet[pageTag] = true
		}
	}

	var tags = make([]string, len(tagsSet))

	i := 0
	for key := range tagsSet {
		tags[i] = key
		i++
	}

	return tags
}

type categoryDto struct {
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
	Pages []string `json:"pages"`
}
