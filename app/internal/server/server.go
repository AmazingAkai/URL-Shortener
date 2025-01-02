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
	"github.com/AmazingAkai/URL-Shortener/app/internal/routes"

	"github.com/go-chi/chi/v5"
	cmiddleware "github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server *http.Server
	db     *sql.DB
}

func New() *Server {
	r := chi.NewRouter()

	r.Use(cmiddleware.RequestID)
	r.Use(cmiddleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(cmiddleware.Recoverer)
	r.Use(middleware.CORS)
	r.Use(middleware.JWT)
	r.Use(middleware.GZip)

	routes.RegisterURLRoutes(r)
	routes.RegisterUserRoutes(r)

	server := &Server{
		server: &http.Server{
			Handler:      r,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		db: database.New(),
	}

	return server
}

func (s *Server) Run() error {
	s.server.Addr = os.Getenv("ADDR")
	log.Infof("Server started at http://%s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %v", err)
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("error closing database: %v", err)
		}
	}

	return nil
}
