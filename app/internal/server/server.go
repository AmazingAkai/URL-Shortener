package server

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
)

type FiberServer struct {
	*fiber.App

	db *sql.DB
}

func New() *FiberServer {
	app := fiber.New(fiber.Config{
		ServerHeader: "github.com/AmazingAkai/URL-Shortener/app",
		AppName:      "github.com/AmazingAkai/URL-Shortener/app",
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	server := &FiberServer{
		App: app,
		db:  database.New(),
	}

	return server
}

func (s *FiberServer) ShutdownWithContext(ctx context.Context) error {
	if err := s.App.ShutdownWithContext(ctx); err != nil {
		return err
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return err
		}
	}

	return nil
}
