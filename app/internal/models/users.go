package models

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=32"`
}

type UserOut struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
