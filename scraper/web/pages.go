package web

import (
	"net/http"

	"github.com/pestanko/miniscrape/scraper"
	"github.com/pestanko/miniscrape/scraper/config"
)

func HandlePages(service *scraper.Service, w http.ResponseWriter, req *http.Request) {
	WriteJsonResponse(w, http.StatusOK, service.Categories)
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
	PageName     string `json:"pageName"`
	PageCodeName string `json:"pageCodeName"`
	HomePage     string `json:"homepage"`
}
