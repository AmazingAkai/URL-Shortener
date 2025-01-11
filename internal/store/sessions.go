package store

import (
	"sync"
)

type Session struct {
	UserID    int
	Token     string
	CSRFToken string
	ExpiresAt int64
}

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
