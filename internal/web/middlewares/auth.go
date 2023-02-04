package middlewares

import (
	"github.com/pestanko/miniscrape/internal/scraper"
	auth2 "github.com/pestanko/miniscrape/internal/web/auth"
	"github.com/pestanko/miniscrape/pkg/rest/webut"
	"net/http"
)

// AuthRequired represents a authentication guard
func AuthRequired(
	_ *scraper.Service,
) func(targetMux http.Handler) http.Handler {
	return func(targetMux http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			session := auth2.GetSessionFromRequest(req)
			if session == nil {
				webut.WriteErrorResponse(w, http.StatusUnauthorized, webut.ErrorDto{
					Error:       "unauthorized",
					ErrorDetail: "You need to login to perform this operation",
				})
				return
			}

			sessionManager := auth2.GetSessionManager()
			if sessionManager.IsSessionValid(*session) {
				targetMux.ServeHTTP(w, req)
				return
			}
		})
	}
}
