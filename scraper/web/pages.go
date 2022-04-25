package web

import (
	"github.com/pestanko/miniscrape/scraper"
	"net/http"
)

func HandlePages(service *scraper.Service, w http.ResponseWriter, req *http.Request) {
	writeJsonResponse(w, http.StatusOK, service.Categories)
}

func HandlePagesContent(service *scraper.Service, w http.ResponseWriter, req *http.Request) {
	selector := scraper.RunSelector{
		Tags:     nil,
		Category: "food",
		Name:     "",
	}

	results := service.Scrape(selector)

	dto := make([]pageContentDto, len(results))

	for i, result := range results {
		dto[i] = pageContentDto{
			Content: result.Content,
			Page: pageContentPageDto{
				PageName:     result.Page.Name,
				PageCodeName: result.Page.CodeName,
				HomePage:     result.Page.Homepage,
			},
			Status: string(result.Status),
		}
	}

	writeJsonResponse(w, http.StatusOK, dto)
}

type pageContentDto struct {
	Content string             `json:"content"`
	Page    pageContentPageDto `json:"page"`
	Status  string             `json:"status"`
}

type pageContentPageDto struct {
	PageName     string `json:"page_name"`
	PageCodeName string `json:"page_code_name"`
	HomePage     string `json:"homepage"`
}
