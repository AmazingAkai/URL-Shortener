package main

import (
	"flag"
	"log"
	"math"

	"github.com/AmazingAkai/URL-Shortener/internal/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := db.New()
	defer db.Close()

	steps := 1
	direction := "down"

	flag.IntVar(&steps, "steps", steps, "number of steps")
	flag.Parse()

	if steps == 0 {
		log.Fatal("steps must not be 0")
	}
	if steps > 0 {
		direction = "up"
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("An error occured initializing the driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("An error occured initializing the migration: %v", err)
	}

	err = m.Steps(steps)
	if err != nil {
		log.Fatalf("An error occured running the migration: %v", err)
	}

	log.Printf("Ran %d migration(s) %s", int(math.Abs(float64(steps))), direction)
}
