package middlewares

import (
	"net/http"
	"time"

	"github.com/pestanko/miniscrape/pkg/config"
	"github.com/pestanko/miniscrape/pkg/utils"
)

const visitorCookie = "VISITOR"

// VisitorCookie middleware
func VisitorCookie(targetMux http.Handler, _ *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		wrapWithVisitorCookie(w, r)

		targetMux.ServeHTTP(w, r)
	})
}

func wrapWithVisitorCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(visitorCookie)

	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:     visitorCookie,
			Value:    utils.RandomString(32),
			Expires:  time.Now().Add(30 * 24 * time.Hour),
			HttpOnly: true,
		}
	}

	http.SetCookie(w, cookie)
}
