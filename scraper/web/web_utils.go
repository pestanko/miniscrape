package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/pestanko/miniscrape/scraper"
	"github.com/pestanko/miniscrape/scraper/utils"
	"github.com/pestanko/miniscrape/scraper/web/auth"
)

type errorDto struct {
	Error       string `json:"error"`
	ErrorDetail string `json:"error_detail"`
}

func WriteErrorResponse(w http.ResponseWriter, code int, err errorDto) {
	log.Printf("Error[%d] - %s: %s", code, err.Error, err.ErrorDetail)

	WriteJsonResponse(w, code, err)
}

func WriteJsonResponse(w http.ResponseWriter, code int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %v", err)
	}
	if _, err := w.Write(jsonResp); err != nil {
		log.Fatalf("Error writing response. Err: %v", err)
	}
}

func WriteUnsupportedHttpMethod(w http.ResponseWriter, method string) {
	WriteErrorResponse(w, http.StatusMethodNotAllowed, errorDto{
		"unsupported_http_method",
		"Unsuppored http method: " + method,
	})
}

func requireAuthentication(service *scraper.Service, w http.ResponseWriter, req *http.Request, callable func()) {
	session := GetSessionFromRequest(req)
	if session != nil {
		sessionManager := auth.GetSessionManager()
		if sessionManager.IsSessionValid(*session) {
			callable()
			return
		}
	}

	WriteErrorResponse(w, http.StatusUnauthorized, errorDto{
		Error:       "unauthorized",
		ErrorDetail: "You need to login to perform this operation",
	})
}

func requireHttpMethod(w http.ResponseWriter, req *http.Request, methods []string, callable func()) {
	if utils.IsInSlice[string](methods, func(i string) bool { return strings.ToUpper(req.Method) == i }) {
		callable()
	} else {
		WriteUnsupportedHttpMethod(w, req.Method)
	}
}
