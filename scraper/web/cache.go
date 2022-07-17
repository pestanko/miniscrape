package web

import (
	"net/http"

	"github.com/pestanko/miniscrape/scraper"
)

// HandleCacheInvalidation handler to handle cache invalidation
func HandleCacheInvalidation(
	service *scraper.Service,
	w http.ResponseWriter,
	req *http.Request,
) {

	selector := makeSelectorFromRequest(req)

	service.InvalidateCache(selector)

	WriteJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "invalidated",
		"message": "cache bas been invalidated",
	})
}
