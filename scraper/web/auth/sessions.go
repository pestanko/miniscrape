package auth

import (
	"time"

	"github.com/pestanko/miniscrape/scraper/utils"
)

var sessionManagerInstance = MakeSessionManagerInstance[SessionData]()

type Session[T any] struct {
	Id         string
	Data       T
	Expiration time.Time
}

func (s Session[T]) IsExpired() bool {
	return s.Expiration.Before(time.Now())
}

type SessionData struct {
	Username string `json:"usernamne"`
}

type SessionManager[T any] interface {
	CreateSession(data T) Session[T]
	GetSession(id string) *Session[T]
	InvalidateSession(id string) bool
	IsSessionValid(session Session[T]) bool
}

type sessionManagerInMemory[T any] struct {
	sessions map[string]Session[T]
}

// CreateSession implements SessionManager
func (s *sessionManagerInMemory[T]) CreateSession(data T) Session[T] {
	now := time.Now()
	exp := now.Add(8 * time.Hour)

	sess := Session[T]{
		Id:         utils.RandomString(64),
		Data:       data,
		Expiration: exp,
	}

	s.sessions[sess.Id] = sess

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
		s.InvalidateSession(session.Id)
		return false
	}

	return true
}

func MakeSessionManagerInstance[T any]() SessionManager[T] {
	return &sessionManagerInMemory[T]{
		sessions: map[string]Session[T]{},
	}
}

func GetSessionManager() SessionManager[SessionData] {
	return sessionManagerInstance
}
