package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pestanko/miniscrape/scraper"
	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/pestanko/miniscrape/scraper/utils"
	"github.com/pestanko/miniscrape/scraper/web/auth"
)

const sessionCookieName = "SESSIONID"

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleAuthLogout(service *scraper.Service, w http.ResponseWriter, req *http.Request) {
	sessionManager := auth.GetSessionManager()
	sessionId := GetSessionIdFromRequest(req)
	if sessionId != "" {
		sessionManager.InvalidateSession(sessionId)
		WriteJsonResponse(w, http.StatusBadRequest, map[string]string{
			"status":  "ok",
			"code":    "error_logout",
			"message": "Unable to logout",
		})
		return
	}

	WriteJsonResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"code":    "success_logout",
		"message": "You have been logged out!",
	})
}

func HandleAuthLogin(service *scraper.Service, w http.ResponseWriter, req *http.Request) {
	users := service.Cfg.Web.Users
	if len(users) == 0 {
		WriteErrorResponse(w, http.StatusBadRequest, errorDto{
			Error:       "unsupported",
			ErrorDetail: "Login is not supported.",
		})
		return
	}

	var cred loginCredentials

	err := json.NewDecoder(req.Body).Decode(&cred)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, errorDto{
			Error:       "invalid_request",
			ErrorDetail: fmt.Sprintf("Error: %v", err),
		})
		return
	}

	user := findUser(users, cred)
	if user == nil {
		WriteErrorResponse(w, http.StatusUnauthorized, errorDto{
			Error:       "unauthorized",
			ErrorDetail: "You have provided invalid credentials",
		})
		return
	}

	sessionManager := auth.GetSessionManager()
	session := sessionManager.CreateSession(auth.SessionData{Username: user.Username})

	sessionCookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    session.Id,
		Domain:   "",
		Expires:  session.Expiration,
		HttpOnly: true,
	}

	http.SetCookie(w, &sessionCookie)
	WriteJsonResponse(w, http.StatusCreated, map[string]string{
		"status":     "created",
		"session_id": session.Id,
	})
}

func HandleSessionStatus(service *scraper.Service, w http.ResponseWriter, req *http.Request) {
	session := GetSessionFromRequest(req)
	if session == nil {
		WriteJsonResponse(w, http.StatusOK, map[string]string{
			"status": "unauthorized",
		})
		return
	}

	sessionManager := auth.GetSessionManager()

	if sessionManager.IsSessionValid(*session) {
		WriteJsonResponse(w, http.StatusOK, map[string]interface{}{
			"status":   "ok",
			"userData": session.Data,
		})
		return
	}

	WriteJsonResponse(w, http.StatusOK, map[string]interface{}{
		"status":  "invalid",
		"expired": session.IsExpired(),
	})
}

func GetSessionIdFromRequest(req *http.Request) string {
	sessCookie := getSessionCookieFromRequest(req)
	if sessCookie == nil {
		return ""
	}

	return sessCookie.Value
}

func GetSessionFromRequest(req *http.Request) *auth.Session[auth.SessionData] {
	sessionManager := auth.GetSessionManager()
	sessionId := GetSessionIdFromRequest(req)
	if sessionId == "" {
		return nil
	}

	return sessionManager.GetSession(sessionId)
}

func getSessionCookieFromRequest(req *http.Request) *http.Cookie {
	sessCookie, err := req.Cookie(sessionCookieName)

	if err != nil && err == http.ErrNoCookie {
		return nil
	}

	return sessCookie
}

func findUser(users []config.User, cred loginCredentials) *config.User {
	return utils.FindInSlice(users, func(u config.User) bool {
		return u.Username == cred.Username && u.Password == u.Password
	})
}
