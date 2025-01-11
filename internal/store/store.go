package store

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users    *UserStore
	Urls     *UrlStore
	Sessions *SessionStore
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:    &UserStore{db},
		Urls:     &UrlStore{db},
		Sessions: &SessionStore{sessions: make(map[string]*Session)},
	}
}
