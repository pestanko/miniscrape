package auth

import (
	"time"

	"github.com/pestanko/miniscrape/scraper/utils"
)

var sessionManagerInstance = MakeSessionManagerInstance[SessionData]()

// Session representation
type Session[T any] struct {
	ID         string
	Data       T
	Expiration time.Time
}

// IsExpired whether the session was expired
func (s Session[T]) IsExpired() bool {
	return s.Expiration.Before(time.Now())
}

// SessionData representation
type SessionData struct {
	Username string `json:"usernamne"`
}

// SessionManager manage the session
type SessionManager[T any] interface {
	// CreateSession with provided session data
	CreateSession(data T) Session[T]
	// GetSession based on the session id
	GetSession(id string) *Session[T]
	// InvalidateSession invalidate session id
	InvalidateSession(id string) bool
	// IsSessionValid check
	IsSessionValid(session Session[T]) bool
}

// MakeSessionManagerInstance creates a new session manager instance
func MakeSessionManagerInstance[T any]() SessionManager[T] {
	return &sessionManagerInMemory[T]{
		sessions: map[string]Session[T]{},
	}
}

// GetSessionManager instance
func GetSessionManager() SessionManager[SessionData] {
	return sessionManagerInstance
}

type sessionManagerInMemory[T any] struct {
	sessions map[string]Session[T]
}

// CreateSession implements SessionManager
func (s *sessionManagerInMemory[T]) CreateSession(data T) Session[T] {
	now := time.Now()
	exp := now.Add(8 * time.Hour)

	sess := Session[T]{
		ID:         utils.RandomString(64),
		Data:       data,
		Expiration: exp,
	}

	s.sessions[sess.ID] = sess

	return sess
}

func (s *sessionManagerInMemory[T]) GetSession(id string) *Session[T] {
	instance, ok := s.sessions[id]
	if !ok {
		return nil
	}
	return &instance
}

func (s *sessionManagerInMemory[T]) InvalidateSession(id string) bool {
	_, ok := s.sessions[id]
	if ok {
		delete(s.sessions, id)
	}

	return ok
}
func (s *sessionManagerInMemory[T]) IsSessionValid(session Session[T]) bool {
	if session.IsExpired() {
		s.InvalidateSession(session.ID)
		return false
	}

	return true
}
