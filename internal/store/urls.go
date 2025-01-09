package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/AmazingAkai/URL-Shortener/internal/utils"
)

type UrlStore struct {
	db *sql.DB
}

type Url struct {
	ID        int        `json:"id"`
	UserID    *int       `json:"user_id,omitempty"`
	LongUrl   string     `json:"long_url"`
	ShortUrl  string     `json:"short_url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type UrlVisit struct {
	UrlID     int    `json:"url_id"`
	IpAddr    string `json:"ip_address"`
	Referer   string `json:"referer"`
	UserAgent string `json:"user_agent"`
}

func (s *UrlStore) Create(ctx context.Context, url *Url) error {
	query := `
		INSERT INTO urls (user_id, long_url, short_url, expires_at) VALUES 
		($1, $2, $3, $4)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	shortUrl, err := s.generateUniqueShortUrl(ctx)
	if err != nil {
		return err
	}

	err = s.db.QueryRowContext(
		ctx,
		query,
		url.UserID,
		url.LongUrl,
		shortUrl,
		url.ExpiresAt,
	).Scan(&url.ID)

	if err != nil {
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "urls_short_url"`:
			return ErrConflict
		default:
			return err
		}
	}

	return nil
}

func (s *UrlStore) GetLongUrl(ctx context.Context, shortUrl string) (*Url, error) {
	query := `
		SELECT id, user_id, long_url, short_url, expires_at FROM urls
		WHERE short_url = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	url := &Url{}
	err := s.db.QueryRowContext(ctx, query, shortUrl).Scan(
		&url.ID,
		&url.UserID,
		&url.LongUrl,
		&url.ShortUrl,
		&url.ExpiresAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return url, nil

}

func (s *UrlStore) CreateVisit(ctx context.Context, visit UrlVisit) {
	query := `INSERT INTO url_requests (url_id, ip_address, referer, user_agent) 
			VALUES ($1, $2, $3, $4)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		visit.UrlID,
		visit.IpAddr,
		visit.IpAddr,
		visit.UserAgent,
	)
	if err != nil {
		log.Printf("Failed to create visit for Url with ID %d: %v", visit.UrlID, err)
	}
}

func (s *UrlStore) generateUniqueShortUrl(ctx context.Context) (string, error) {
	var (
		shortUrl string
		attempts = 0
	)

	for {
		if attempts >= 10 {
			return "", fmt.Errorf("failed to generate unique short Url after %d attempts", attempts)
		}

		shortUrl = utils.GenerateShortUrl()
		_, err := s.GetLongUrl(ctx, shortUrl)
		if err != nil {
			if err == ErrNotFound {
				break
			} else {
				return "", err
			}
		}

		attempts++
	}

	return shortUrl, nil
}
