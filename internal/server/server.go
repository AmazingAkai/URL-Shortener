package server

import (
	"context"
	"database/sql"
	"fmt"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/AmazingAkai/URL-Shortener/internal/db"
	"github.com/AmazingAkai/URL-Shortener/internal/middleware"
	"github.com/AmazingAkai/URL-Shortener/internal/store"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	*http.Server
	db    *sql.DB
	store *store.Storage
}

func New() *Server {
	r := chi.NewRouter()

	db := db.New()
	storage := store.NewStorage(db)

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Use(middleware.CORS)
	r.Use(middleware.Auth(storage))

	server := &Server{
		Server: &http.Server{
			Handler:      r,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		db:    db,
		store: storage,
	}

	// Pages
	r.Get("/", server.indexHandler)
	r.Get("/register", server.registerPageHandler)
	r.Get("/login", server.loginPageHandler)

	// API
	r.Get("/{short_url}", server.redirectShortUrlHandler)
	r.Post("/urls", server.createShortUrlHandler)
	r.Get("/urls", server.getShortUrlList)
	r.Delete("/urls/{url_id}", server.deleteShortUrl)

	r.Post("/register", server.createUserHandler)
	r.Post("/login", server.loginHandler)
	r.Get("/logout", server.logoutHandler)

	// Static
	r.Mount("/static", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	r.NotFound(server.notFoundHandler)

	return server
}

func (s *Server) Run() error {
	s.Addr = os.Getenv("ADDR")
	log.Printf("Server started at http://%s", s.Addr)
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
