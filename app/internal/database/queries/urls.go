package queries

import (
	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
)

func CreateShortURL(url models.URL) error {
	_, err := database.New().Exec("INSERT INTO urls (original_url, short_url) VALUES ($1, $2)", url.LongURL, url.ShortURL)
	return err
}

func GetLongURL(shortURL string) (string, error) {
	var longURL string
	err := database.New().QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&longURL)
	return longURL, err
}

func ShortURLExists(shortURL string) (bool, error) {
	var exists bool
	err := database.New().QueryRow("SELECT EXISTS (SELECT 1 FROM urls WHERE short_url = $1)", shortURL).Scan(&exists)
	return exists, err
}
