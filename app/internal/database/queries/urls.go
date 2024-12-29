package queries

import (
	"database/sql"
	"time"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
)

func CreateShortURL(url *models.URL) error {
	query := `
		INSERT INTO urls (original_url, short_url, expires_at) 
		VALUES ($1, $2, $3) ON CONFLICT (short_url) 
		DO UPDATE SET original_url = EXCLUDED.original_url, expires_at = EXCLUDED.expires_at
		RETURNING id, original_url, short_url, created_at, expires_at`

	err := database.New().QueryRow(query, url.LongURL, url.ShortURL, url.ExpiresAt).Scan(
		&url.ID, &url.LongURL, &url.ShortURL, &url.CreatedAt, &url.ExpiresAt,
	)
	if err != nil {
		return err
	}
	return nil
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
