package handlers

import (
	"net/http"

	"github.com/pestanko/miniscrape/internal/models"
	"github.com/pestanko/miniscrape/internal/scraper"
	"github.com/pestanko/miniscrape/pkg/rest/webut"
)

// HandlePages handler
func HandlePages(service *scraper.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webut.WriteJSONResponse(w, http.StatusOK, service.GetCategories())
	}
}

// HandlePagesContent handler
func HandlePagesContent(service *scraper.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		selector := makeSelectorFromRequest(req)

		results := service.Scrape(req.Context(), selector)

		dto := make([]pageContentDto, len(results))

		for i, result := range results {
			dto[i] = pageContentDto{
				Content:  result.Content,
				Status:   string(result.Status),
				Resolver: result.Page.Resolver,
				Page: pageContentPageDto{
					PageName:     result.Page.Name,
					PageCodeName: result.Page.CodeName,
					URL:          result.Page.URL,
					HomePage:     result.Page.Homepage,
					Tags:         result.Page.Tags,
					Category:     result.Page.Category,
				},
			}
		}

		webut.WriteJSONResponse(w, http.StatusOK, dto)
	}
}

type pageContentDto struct {
	Content  string             `json:"content"`
	Status   string             `json:"status"`
	Resolver string             `json:"resolver"`
	Page     pageContentPageDto `json:"page"`
}

type pageContentPageDto struct {
	PageName     string   `json:"name"`
	PageCodeName string   `json:"codename"`
	HomePage     string   `json:"homepage"`
	URL          string   `json:"url"`
	Tags         []string `json:"tags"`
	Category     string   `json:"category"`
}

func makeSelectorFromRequest(req *http.Request) models.RunSelector {
	category := req.URL.Query().Get("c")
	tags := req.URL.Query()["t"]
	name := req.URL.Query().Get("n")

	return models.RunSelector{
		Tags:     tags,
		Category: category,
		Page:     name,
	}
}
