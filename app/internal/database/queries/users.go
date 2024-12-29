package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
)

func CreateUser(user *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, password, created_at`

	err = database.New().QueryRow(query, user.Email, hash).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
			return errors.New("email already exists")
		}
		return err
	}

	return nil
}

func AuthenticateUser(user *models.User) error {
	var hashedPassword string
	err := database.New().QueryRow("SELECT password FROM users WHERE email = $1", user.Email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return errors.New("incorrect password")
	}

	return nil
}
