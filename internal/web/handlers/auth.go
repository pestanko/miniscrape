package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pestanko/miniscrape/internal/scraper"
	auth2 "github.com/pestanko/miniscrape/internal/web/auth"
	"github.com/pestanko/miniscrape/pkg/rest/webut"
	"net/http"
)

// HandleAuthLogout logout handler
func HandleAuthLogout(_ *scraper.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sessionManager := auth2.GetSessionManager()
		sessionID := auth2.GetSessionIDFromRequest(req)
		if sessionID != "" {
			sessionManager.InvalidateSession(sessionID)
			webut.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{
				"status":  "ok",
				"code":    "error_logout",
				"message": "Unable to logout",
			})
			return
		}

		webut.WriteJSONResponse(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"code":    "success_logout",
			"message": "You have been logged out!",
		})
	}
}

// HandleAuthLogin login handler
func HandleAuthLogin(service *scraper.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		users := service.Cfg.Web.Users
		if len(users) == 0 {
			webut.WriteErrorResponse(w, http.StatusBadRequest, webut.ErrorDto{
				Error:       "unsupported",
				ErrorDetail: "Login is not supported.",
			})
			return
		}

		var cred auth2.LoginCredentials

		err := json.NewDecoder(req.Body).Decode(&cred)
		if err != nil {
			webut.WriteErrorResponse(w, http.StatusBadRequest, webut.ErrorDto{
				Error:       "invalid_request",
				ErrorDetail: fmt.Sprintf("Error: %v", err),
			})
			return
		}

		user := auth2.FindUser(users, cred)
		if user == nil {
			webut.WriteErrorResponse(w, http.StatusUnauthorized, webut.ErrorDto{
				Error:       "unauthorized",
				ErrorDetail: "You have provided invalid credentials",
			})
			return
		}

		sessionManager := auth2.GetSessionManager()
		session := sessionManager.CreateSession(auth2.SessionData{Username: user.Username})

		sessionCookie := auth2.CreateSessionCookie(session)

		http.SetCookie(w, &sessionCookie)
		webut.WriteJSONResponse(w, http.StatusCreated, map[string]string{
			"status":     "created",
			"session_id": session.ID,
		})
	}
}

// HandleSessionStatus session status handler
func HandleSessionStatus(_ *scraper.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session := auth2.GetSessionFromRequest(req)
		if session == nil {
			webut.WriteJSONResponse(w, http.StatusOK, map[string]string{
				"status": "unauthorized",
			})
			return
		}

		sessionManager := auth2.GetSessionManager()

		if sessionManager.IsSessionValid(*session) {
			webut.WriteJSONResponse(w, http.StatusOK, map[string]any{
				"status":   "ok",
				"userData": session.Data,
			})
			return
		}

		webut.WriteJSONResponse(w, http.StatusOK, map[string]any{
			"status":  "invalid",
			"expired": session.IsExpired(),
		})
	}
}
