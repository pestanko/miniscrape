package handlers

import (
	"net/http"

	"github.com/pestanko/miniscrape/pkg"
	"github.com/pestanko/miniscrape/pkg/web/wutt"
)

// HandleCacheInvalidation handler to handle cache invalidation
func HandleCacheInvalidation(service *pkg.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		selector := makeSelectorFromRequest(req)

		service.InvalidateCache(selector)

		wutt.WriteJSONResponse(w, http.StatusOK, map[string]string{
			"status":  "invalidated",
			"message": "cache bas been invalidated",
		})
	}
}
