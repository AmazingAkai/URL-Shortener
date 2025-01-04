package server

import (
	"context"
	"database/sql"
	"fmt"

	"net/http"
	"os"
	"time"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
	"github.com/AmazingAkai/URL-Shortener/app/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	*http.Server
	db *sql.DB
}

func New() *Server {
	r := chi.NewRouter()
	db := database.New()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORS)
	r.Use(middleware.JWT)
	r.Use(middleware.GZip)

	server := &Server{
		Server: &http.Server{
			Handler:      r,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		db: db,
	}

	server.RegisterURLRoutes(r)
	server.RegisterUserRoutes(r)

	return server
}

func (s *Server) Run() error {
	s.Addr = os.Getenv("ADDR")
	log.Infof("Server started at http://%s", s.Addr)
	return s.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %v", err)
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("error closing database: %v", err)
		}
	}

	return nil
}
