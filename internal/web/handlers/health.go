package handlers

import (
	"github.com/pestanko/miniscrape/pkg/rest/webut"
	"net/http"
)

// HandleHealthStatus handler
func HandleHealthStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		resp := make(map[string]string)
		resp["status"] = "ok"

		webut.WriteJSONResponse(w, http.StatusOK, resp)
	}
}
