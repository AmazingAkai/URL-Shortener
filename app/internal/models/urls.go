package models

import "time"

type URL struct {
	ID        int        `json:"id"`
	ShortURL  string     `json:"short_url" validate:"required,min=8,max=8"`
	LongURL   string     `json:"long_url" validate:"required,url"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"` // TODO: If given, should be future date + valid date
}
