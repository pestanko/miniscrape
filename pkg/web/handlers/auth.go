package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pestanko/miniscrape/internal/scraper"
	"net/http"

	"github.com/pestanko/miniscrape/pkg/web/auth"
	"github.com/pestanko/miniscrape/pkg/web/wutt"
)

// HandleAuthLogout logout handler
func HandleAuthLogout(_ *scraper.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sessionManager := auth.GetSessionManager()
		sessionID := auth.GetSessionIDFromRequest(req)
		if sessionID != "" {
			sessionManager.InvalidateSession(sessionID)
			wutt.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{
				"status":  "ok",
				"code":    "error_logout",
				"message": "Unable to logout",
			})
			return
		}

		wutt.WriteJSONResponse(w, http.StatusOK, map[string]string{
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
			wutt.WriteErrorResponse(w, http.StatusBadRequest, wutt.ErrorDto{
				Error:       "unsupported",
				ErrorDetail: "Login is not supported.",
			})
			return
		}

		var cred auth.LoginCredentials

		err := json.NewDecoder(req.Body).Decode(&cred)
		if err != nil {
			wutt.WriteErrorResponse(w, http.StatusBadRequest, wutt.ErrorDto{
				Error:       "invalid_request",
				ErrorDetail: fmt.Sprintf("Error: %v", err),
			})
			return
		}

		user := auth.FindUser(users, cred)
		if user == nil {
			wutt.WriteErrorResponse(w, http.StatusUnauthorized, wutt.ErrorDto{
				Error:       "unauthorized",
				ErrorDetail: "You have provided invalid credentials",
			})
			return
		}

		sessionManager := auth.GetSessionManager()
		session := sessionManager.CreateSession(auth.SessionData{Username: user.Username})

		sessionCookie := auth.CreateSessionCookie(session)

		http.SetCookie(w, &sessionCookie)
		wutt.WriteJSONResponse(w, http.StatusCreated, map[string]string{
			"status":     "created",
			"session_id": session.ID,
		})
	}
}

// HandleSessionStatus session status handler
func HandleSessionStatus(_ *scraper.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session := auth.GetSessionFromRequest(req)
		if session == nil {
			wutt.WriteJSONResponse(w, http.StatusOK, map[string]string{
				"status": "unauthorized",
			})
			return
		}

		sessionManager := auth.GetSessionManager()

		if sessionManager.IsSessionValid(*session) {
			wutt.WriteJSONResponse(w, http.StatusOK, map[string]any{
				"status":   "ok",
				"userData": session.Data,
			})
			return
		}

		wutt.WriteJSONResponse(w, http.StatusOK, map[string]any{
			"status":  "invalid",
			"expired": session.IsExpired(),
		})
	}
}
