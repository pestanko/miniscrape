package middlewares

import (
	"net/http"

	"github.com/pestanko/miniscrape/scraper/config"
)

const (
	headerAccessControlAllowOrigin  = "Access-Control-Allow-Origin"
	headerAccessControlAllowMethods = "Access-Control-Allow-Methods"
	headerAccessControlAllowHeaders = "Access-Control-Allow-Headers"

	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	allowedOrigin  = "*"
	allowedMethods = "POST, GET, OPTIONS, PUT, DELETE"
)

// SetupCors middleware
func SetupCors(targetMux http.Handler, cfg *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(headerAccessControlAllowOrigin, allowedOrigin)
		w.Header().Set(headerAccessControlAllowMethods, allowedMethods)
		w.Header().Set(headerAccessControlAllowHeaders, allowedHeaders)

		targetMux.ServeHTTP(w, req)
	})
}
