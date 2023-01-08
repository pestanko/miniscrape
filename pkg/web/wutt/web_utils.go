package wutt

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// ErrorDto represents an error response
type ErrorDto struct {
	Error       string `json:"error"`
	ErrorDetail string `json:"error_detail"`
}

// WriteErrorResponse helper
func WriteErrorResponse(w http.ResponseWriter, code int, err ErrorDto) {
	log.Warn().
		Str("error", err.Error).
		Str("detail", err.ErrorDetail).
		Int("code", code).
		Msg("Returning the error response")

	WriteJSONResponse(w, code, err)
}

// WriteJSONResponse helper
func WriteJSONResponse(w http.ResponseWriter, code int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Error().Err(err).Msg("Error happened in JSON marshal")
	}
	if _, err := w.Write(jsonResp); err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}
}
