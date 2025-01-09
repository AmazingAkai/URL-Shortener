package store

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int      `json:"id"`
	Email    string   `json:"email"`
	Password password `json:"-"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password) VALUES 
		($1, $2)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password.hash,
	).Scan(&user.ID)

	if err != nil {
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrConflict
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password.hash,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
