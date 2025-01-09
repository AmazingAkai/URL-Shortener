package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AmazingAkai/URL-Shortener/internal/server"

	_ "github.com/joho/godotenv/autoload"
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
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Printf("Server stopped")
}
