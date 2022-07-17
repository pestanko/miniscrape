package middlewares

import (
	"net/http"

	"github.com/pestanko/miniscrape/pkg/config"
)

// Middleware main middleware representation
type Middleware func(handler http.Handler, cfg *config.AppConfig) http.Handler

// ApplyMiddlewares all the middlewares
func ApplyMiddlewares(targetMux http.Handler, cfg *config.AppConfig, middlewares []Middleware) http.Handler {
	result := targetMux
	for _, middleware := range middlewares {
		result = middleware(result, cfg)
	}

	return result
}
