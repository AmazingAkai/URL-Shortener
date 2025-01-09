package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	db *sql.DB
)

func New() *sql.DB {
	if db != nil {
		return db
	}

	var err error
	db, err = sql.Open("pgx", os.Getenv("DATABASE_URI"))
	if err != nil {
		log.Fatalf("Failed to open the database connection: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return db
}

func Close() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close the database connection: %v", err)
		}
	}
}
