package web

import (
	"github.com/pestanko/miniscrape/scraper"
	"net/http"
)

func HandlePages(service *scraper.Service, w http.ResponseWriter, req *http.Request) {
	writeJsonResponse(w, http.StatusOK, service.Categories)
}
