package web

import (
	"net/http"

	"github.com/pestanko/miniscrape/pkg"
)

// HandleCacheInvalidation handler to handle cache invalidation
func HandleCacheInvalidation(
	service *pkg.Service,
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
