package models

import "time"

type URL struct {
	LongURL   string     `json:"long_url" validate:"required,url"`
	ExpiresAt *time.Time `json:"expires_at" validate:"omitempty,futureDate"`
}

type URLOut struct {
	ID        int        `json:"id"`
	UserID    *int       `json:"user_id,omitempty"`
	ShortURL  string     `json:"short_url"`
	LongURL   string     `json:"long_url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}
