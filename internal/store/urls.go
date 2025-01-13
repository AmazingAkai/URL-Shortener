package store

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

type UrlStore struct {
	db *sql.DB
}

type Url struct {
	ID       int
	UserID   *int
	LongUrl  string
	ShortUrl string
	Visits   int
}

func (s *UrlStore) Create(ctx context.Context, url *Url) error {
	query := `
		INSERT INTO urls (user_id, long_url, short_url) VALUES 
		($1, $2, $3)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		url.UserID,
		url.LongUrl,
		url.ShortUrl,
	).Scan(&url.ID)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value"):
			return ErrConflict
		default:
			return err
		}
	}

	return nil
}

func (s *UrlStore) IncrementVisits(urlID int) {
	query := "UPDATE urls SET visits = visits + 1 WHERE id = $1"

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, urlID)
	if err != nil {
		log.Printf("Failed to increment visits for url %d: %v", urlID, err)
	}
}

func (s *UrlStore) GetUrl(ctx context.Context, shortUrl string) (*Url, error) {
	query := `
		SELECT id, user_id, long_url, short_url, visits FROM urls
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
		&url.Visits,
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

func (s *UrlStore) GetUrlList(ctx context.Context, userID int) ([]*Url, error) {
	query := `
		SELECT id, user_id, long_url, short_url, visits FROM urls
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*Url
	for rows.Next() {
		url := &Url{}
		err := rows.Scan(
			&url.ID,
			&url.UserID,
			&url.LongUrl,
			&url.ShortUrl,
			&url.Visits,
		)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
func (s *UrlStore) Delete(ctx context.Context, urlID, userID int) error {
	query := `DELETE FROM urls WHERE id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, urlID, userID)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return ErrNotFound
	}

	return nil
}
