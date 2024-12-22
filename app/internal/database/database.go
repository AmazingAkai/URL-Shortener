package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	dbUri      = os.Getenv("DATABASE_URI")
	dbInstance *sql.DB
)

func New() *sql.DB {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	dbInstance, err := sql.Open("pgx", dbUri)
	if err != nil {
		log.Fatal(err)
	}

	return dbInstance
}
