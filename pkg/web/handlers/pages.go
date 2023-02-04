package handlers

import (
	"github.com/pestanko/miniscrape/internal/config"
	"net/http"

	"github.com/pestanko/miniscrape/pkg"
	"github.com/pestanko/miniscrape/pkg/web/wutt"
)

// HandlePages handler
func HandlePages(service *pkg.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wutt.WriteJSONResponse(w, http.StatusOK, service.GetCategories())
	}
}

// HandlePagesContent handler
func HandlePagesContent(service *pkg.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		selector := makeSelectorFromRequest(req)

		results := service.Scrape(req.Context(), selector)

		dto := make([]pageContentDto, len(results))

		for i, result := range results {
			dto[i] = pageContentDto{
				Content: result.Content,
				Status:  string(result.Status),
				Page: pageContentPageDto{
					PageName:     result.Page.Name,
					PageCodeName: result.Page.CodeName,
					HomePage:     result.Page.Homepage,
					Tags:         result.Page.Tags,
					Category:     result.Page.Category,
				},
			}
		}

		wutt.WriteJSONResponse(w, http.StatusOK, dto)
	}
}

type pageContentDto struct {
	Content string             `json:"content"`
	Status  string             `json:"status"`
	Page    pageContentPageDto `json:"page"`
}

type pageContentPageDto struct {
	PageName     string   `json:"name"`
	PageCodeName string   `json:"codename"`
	HomePage     string   `json:"homepage"`
	Tags         []string `json:"tags"`
	Category     string   `json:"category"`
}

func makeSelectorFromRequest(req *http.Request) config.RunSelector {
	category := req.URL.Query().Get("c")
	tags := req.URL.Query()["t"]
	name := req.URL.Query().Get("n")

	return config.RunSelector{
		Tags:     tags,
		Category: category,
		Page:     name,
	}
}
