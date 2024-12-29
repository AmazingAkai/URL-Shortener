package database

import (
	"database/sql"

	"os"

	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	db *sql.DB
)

func New() *sql.DB {
	if db != nil {
		return db
	}

	DB, err := sql.Open("pgx", os.Getenv("DATABASE_URI"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return DB
}
