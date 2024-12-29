package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required,min=8,max=32"`
	CreatedAt time.Time `json:"created_at"`
}
