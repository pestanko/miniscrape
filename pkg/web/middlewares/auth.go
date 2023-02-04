package middlewares

import (
	"github.com/pestanko/miniscrape/internal/scraper"
	"net/http"

	"github.com/pestanko/miniscrape/pkg/web/auth"
	"github.com/pestanko/miniscrape/pkg/web/wutt"
)

// AuthRequired represents a authentication guard
func AuthRequired(
	_ *scraper.Service,
) func(targetMux http.Handler) http.Handler {
	return func(targetMux http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			session := auth.GetSessionFromRequest(req)
			if session == nil {
				wutt.WriteErrorResponse(w, http.StatusUnauthorized, wutt.ErrorDto{
					Error:       "unauthorized",
					ErrorDetail: "You need to login to perform this operation",
				})
				return
			}

			sessionManager := auth.GetSessionManager()
			if sessionManager.IsSessionValid(*session) {
				targetMux.ServeHTTP(w, req)
				return
			}
		})
	}
}
