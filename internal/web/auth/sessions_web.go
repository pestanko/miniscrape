package auth

import "net/http"

const sessionCookieName = "SESSIONID"

// GetSessionIDFromRequest get session id
func GetSessionIDFromRequest(req *http.Request) string {
	sessCookie := getSessionCookieFromRequest(req)
	if sessCookie == nil {
		return ""
	}

	return sessCookie.Value
}

// GetSessionFromRequest get whole session with auth data
func GetSessionFromRequest(req *http.Request) *Session[SessionData] {
	sessionManager := GetSessionManager()
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

// CreateSessionCookie create session cookie with the session data
func CreateSessionCookie(session Session[SessionData]) http.Cookie {
	sessionCookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    session.ID,
		Domain:   "",
		Expires:  session.Expiration,
		HttpOnly: true,
	}
	return sessionCookie
}
