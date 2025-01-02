package queries

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"

	"golang.org/x/crypto/bcrypt"
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
		RETURNING id, email`

	err = database.New().
		QueryRow(query, userInput.Email, hash).
		Scan(&user.ID, &user.Email)

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
		QueryRow("SELECT id, email, password FROM users WHERE email = $1", userIn.Email).
		Scan(&user.ID, &user.Email, &hashedPassword)

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

	return user, nil
}
