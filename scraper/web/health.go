package web

import (
	"net/http"
)

func HandleHealthStatus(w http.ResponseWriter, request *http.Request) {
	resp := make(map[string]string)
	resp["status"] = "active"

	writeJsonResponse(w, 200, resp)
}
