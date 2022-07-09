package middlewares

import (
	"net/http"

	"github.com/pestanko/miniscrape/scraper/config"
)

func SetupCors(targetMux http.Handler, cfg *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		targetMux.ServeHTTP(w, req)
	})
}
