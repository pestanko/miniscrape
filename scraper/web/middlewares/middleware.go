package middlewares

import (
	"net/http"

	"github.com/pestanko/miniscrape/scraper/config"
)

type Middleware func(handler http.Handler, cfg *config.AppConfig) http.Handler

func ApplyMiddlewares(targetMux http.Handler, cfg *config.AppConfig, middlewares []Middleware) http.Handler {
	result := targetMux
	for _, middleware := range middlewares {
		result = middleware(result, cfg)
	}

	return result
}
