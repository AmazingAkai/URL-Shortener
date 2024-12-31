package queries

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils"
)

func CreateShortURL(urlInput models.URL, user any) (*models.URLOut, error) {
	var (
		userID *int = nil
		url         = &models.URLOut{}
	)
	if user != nil {
		userID = &user.(*models.UserOut).ID
	}

	shortURL, err := generateUniqueShortURL()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO urls (user_id, original_url, short_url, expires_at) 
		VALUES ($1, $2, $3, $4) ON CONFLICT (short_url) 
		DO UPDATE SET original_url = EXCLUDED.original_url, expires_at = EXCLUDED.expires_at
		RETURNING id, user_id, original_url, short_url, expires_at`

	err = database.New().
		QueryRow(query, userID, urlInput.LongURL, shortURL, urlInput.ExpiresAt).
		Scan(&url.ID, &url.UserID, &url.LongURL, &url.ShortURL, &url.ExpiresAt)

	if err != nil {
		return nil, err
	}
	return url, nil
}

func GetLongURL(shortURL string) (string, error) {
	var (
		longURL   string
		expiresAt *time.Time
	)
	query := "SELECT original_url, expires_at FROM urls WHERE short_url = $1 LIMIT 1"
	err := database.New().QueryRow(query, shortURL).Scan(&longURL, &expiresAt)
	if err == sql.ErrNoRows {
		return "", nil
	}

	if expiresAt != nil && time.Now().After(*expiresAt) {
		return "", nil
	}

	return longURL, err
}

func generateUniqueShortURL() (string, error) {
	var (
		shortURL string
		attempts = 0
	)

	for {
		if attempts >= 10 {
			return "", fmt.Errorf("failed to generate unique short URL after %d attempts", attempts)
		}

		shortURL = utils.GenerateShortURL()
		longURL, err := GetLongURL(shortURL)
		if err != nil {
			return "", err
		}

		if longURL == "" {
			break
		}
		attempts++
	}

	return shortURL, nil
}
