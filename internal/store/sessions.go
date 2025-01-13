package store

import (
	"sync"
	"time"
)

type Session struct {
	UserID    int
	Token     string
	ExpiresAt int64
}

// TODO: Autremove expired sessions
type SessionStore struct {
	sessions   map[string]*Session
	sessionsMu sync.RWMutex
}

func (s *SessionStore) Get(token string) (*Session, error) {
	s.sessionsMu.RLock()
	defer s.sessionsMu.RUnlock()

	session, ok := s.sessions[token]
	if !ok {
		return nil, ErrNotFound
	}

	if session.ExpiresAt < time.Now().Unix() {
		s.Delete(token)
		return nil, ErrNotFound
	}

	return session, nil
}

func (s *SessionStore) Set(token string, session *Session) {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	s.sessions[token] = session
}

func (s *SessionStore) Delete(token string) {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	delete(s.sessions, token)
}
