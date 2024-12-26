package main

import (
	"flag"
	"log"
	"math"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/AmazingAkai/URL-Shortener/app/internal/autoload"
	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
)

func main() {
	db := database.New()
	defer db.Close()

	steps := flag.Int("steps", 1, "number of steps")
	flag.Parse()

	direction := "down"
	if *steps > 0 {
		direction = "up"
	} else if *steps == 0 {
		log.Fatal("steps must not be 0")
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

	err = m.Steps(*steps)
	if err != nil {
		log.Fatalf("An error occured running the migration: %v", err)
	}

	log.Printf("Ran %d migration(s) %s", int(math.Abs(float64(*steps))), direction)
}
