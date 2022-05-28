package web

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorDto struct {
	Error       string `json:"error"`
	ErrorDetail string `json:"error_detail"`
}

func WriteErrorResponse(w http.ResponseWriter, code int, error errorDto) {
	log.Printf("Error[%d] - %s: %s", code, error.Error, error.ErrorDetail)

	WriteJsonResponse(w, code, error)
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
