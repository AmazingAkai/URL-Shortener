package main

import (
	"context"
	"net/http"
	"os"

	"os/signal"
	"syscall"
	"time"

	_ "github.com/AmazingAkai/URL-Shortener/app/internal/env"
	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
	"github.com/AmazingAkai/URL-Shortener/app/internal/server"
)

func main() {
	s := server.New()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error running server: %v", err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Infof("Server forced to shutdown with error: %v", err)
	}

	log.Info("Server stopped")
}
