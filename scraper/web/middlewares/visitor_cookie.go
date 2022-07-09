package middlewares

import (
	"net/http"
	"time"

	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/pestanko/miniscrape/scraper/utils"
)

const visitorCookie = "VISITOR"

func VisitorCookie(targetMux http.Handler, cfg *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		WrapWithVisitorCookie(w, r)

		targetMux.ServeHTTP(w, r)
	})
}

func WrapWithVisitorCookie(w http.ResponseWriter, r *http.Request) {
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
