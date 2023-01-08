package handlers

import (
	"net/http"

	"github.com/pestanko/miniscrape/pkg/web/wutt"
)

// HandleHealthStatus handler
func HandleHealthStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		resp := make(map[string]string)
		resp["status"] = "ok"

		wutt.WriteJSONResponse(w, http.StatusOK, resp)
	}
}
