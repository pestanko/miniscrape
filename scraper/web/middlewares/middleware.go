package middlewares

import (
	"net/http"
)

type Middleware func(handler http.Handler) http.Handler

func ApplyMiddlewares(targetMux http.Handler, middlewares []Middleware) http.Handler {
	result := targetMux
	for _, middleware := range middlewares {
		result = middleware(result)
	}

	return result
}
