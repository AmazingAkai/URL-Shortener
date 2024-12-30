package queries

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
)

func CreateUser(userInput models.User) (*models.UserOut, error) {
	user := &models.UserOut{}
	hash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, created_at`

	err = database.New().
		QueryRow(query, userInput.Email, hash).
		Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func AuthenticateUser(userIn models.User) (*models.UserOut, error) {
	var (
		user           = &models.UserOut{}
		hashedPassword string
	)

	err := database.New().
		QueryRow("SELECT id, email, password, created_at FROM users WHERE email = $1", userIn.Email).
		Scan(&user.ID, &user.Email, &hashedPassword, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userIn.Password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	return nil, nil
}