package web

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func writeJsonResponse(w http.ResponseWriter, code int, resp interface{}) {
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

func applyMiddlewares(targetMux http.Handler) http.Handler {
	return RequestLogger(targetMux)
}

func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		targetMux.ServeHTTP(w, r)

		// log request by who(IP address)
		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t\t%s\t\t%s\t\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
	})
}
