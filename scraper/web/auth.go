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

// HandleAuthLogout logout handler
func HandleAuthLogout(_ *scraper.Service, w http.ResponseWriter, req *http.Request) {
	sessionManager := auth.GetSessionManager()
	sessionID := GetSessionIDFromRequest(req)
	if sessionID != "" {
		sessionManager.InvalidateSession(sessionID)
		WriteJSONResponse(w, http.StatusBadRequest, map[string]string{
			"status":  "ok",
			"code":    "error_logout",
			"message": "Unable to logout",
		})
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"code":    "success_logout",
		"message": "You have been logged out!",
	})
}

// HandleAuthLogin login handler
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
		Value:    session.ID,
		Domain:   "",
		Expires:  session.Expiration,
		HttpOnly: true,
	}

	http.SetCookie(w, &sessionCookie)
	WriteJSONResponse(w, http.StatusCreated, map[string]string{
		"status":     "created",
		"session_id": session.ID,
	})
}

// HandleSessionStatus session status handler
func HandleSessionStatus(_ *scraper.Service, w http.ResponseWriter, req *http.Request) {
	session := GetSessionFromRequest(req)
	if session == nil {
		WriteJSONResponse(w, http.StatusOK, map[string]string{
			"status": "unauthorized",
		})
		return
	}

	sessionManager := auth.GetSessionManager()

	if sessionManager.IsSessionValid(*session) {
		WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
			"status":   "ok",
			"userData": session.Data,
		})
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"status":  "invalid",
		"expired": session.IsExpired(),
	})
}

// GetSessionIDFromRequest get session id
func GetSessionIDFromRequest(req *http.Request) string {
	sessCookie := getSessionCookieFromRequest(req)
	if sessCookie == nil {
		return ""
	}

	return sessCookie.Value
}

// GetSessionFromRequest get whole session with auth data
func GetSessionFromRequest(req *http.Request) *auth.Session[auth.SessionData] {
	sessionManager := auth.GetSessionManager()
	sessionID := GetSessionIDFromRequest(req)
	if sessionID == "" {
		return nil
	}

	return sessionManager.GetSession(sessionID)
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
		return u.Username == cred.Username && u.Password == cred.Password
	})
}
