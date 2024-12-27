package main

import (
	"flag"

	"math"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils"
)

func main() {
	utils.LoadEnv()

	db := database.New()
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

	log.Infof("Ran %d migration(s) %s", int(math.Abs(float64(steps))), direction)
}
