package web

import (
	"net/http"

	"github.com/pestanko/miniscrape/scraper"
	"github.com/pestanko/miniscrape/scraper/config"
)

func HandlePages(service *scraper.Service, w http.ResponseWriter, req *http.Request) {
	WriteJsonResponse(w, http.StatusOK, service.GetCategories())
}

func HandlePagesContent(service *scraper.Service, w http.ResponseWriter, req *http.Request) {
	category := req.URL.Query().Get("c")
	tags := req.URL.Query()["t"]
	name := req.URL.Query().Get("n")

	selector := config.RunSelector{
		Tags:     tags,
		Category: category,
		Page:     name,
	}

	results := service.Scrape(selector)

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

	WriteJsonResponse(w, http.StatusOK, dto)
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
