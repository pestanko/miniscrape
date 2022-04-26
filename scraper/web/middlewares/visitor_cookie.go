package middlewares

import (
	"github.com/pestanko/miniscrape/scraper/utils"
	"net/http"
	"time"
)

const visitorCookie = "VISITOR"

func VisitorCookie(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		targetMux.ServeHTTP(w, r)
	})
}
