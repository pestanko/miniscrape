package handlers

import (
	"net/http"
	"time"

	"github.com/pestanko/miniscrape/pkg/rest/webut"
)

var startTime time.Time

// HandleHealthStatus handler
func HandleHealthStatus() http.HandlerFunc {
	if startTime.IsZero() {
		startTime = time.Now()
	}

	return func(w http.ResponseWriter, _ *http.Request) {
		resp := make(map[string]string)
		resp["status"] = "ok"
		resp["uptime"] = time.Since(startTime).String()
		resp["started_at"] = startTime.Format(time.RFC3339)

		webut.WriteJSONResponse(w, http.StatusOK, resp)
	}
}
