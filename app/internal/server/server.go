package server

import (
	"context"
	"database/sql"
	"fmt"

	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database"
	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
	"github.com/AmazingAkai/URL-Shortener/app/internal/middleware"
	"github.com/AmazingAkai/URL-Shortener/app/internal/routes"

	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	db     *sql.DB
}

func New() *Server {
	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggerMiddleware)
	router.Use(middleware.JwtMiddleware)
	router.Use(middleware.GzipMiddleware)

	routes.RegisterURLRoutes(router)
	routes.RegisterUserRoutes(router)

	server := &Server{
		server: &http.Server{
			Handler:      router,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		db: database.New(),
	}

	return server
}

func (s *Server) Run() error {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	s.server.Addr = fmt.Sprintf(":%d", port)

	log.Infof("Server starting on port %d", port)
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
