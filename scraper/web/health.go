package web

import (
	"net/http"
)

// HandleHealthStatus handler
func HandleHealthStatus(w http.ResponseWriter, _ *http.Request) {
	resp := make(map[string]string)
	resp["status"] = "active"

	WriteJSONResponse(w, 200, resp)
}
