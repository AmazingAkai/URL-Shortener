package main

import (
	"context"
	"fmt"

	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
	"github.com/AmazingAkai/URL-Shortener/app/internal/server"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils"
)

func gracefulShutdown(s *server.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Info("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Infof("Server forced to shutdown with error: %v", err)
	}

	log.Info("Server exiting")

	done <- true
}

func main() {
	utils.LoadEnv()

	server := server.New()
	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	err := server.Run()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Infof("Graceful shutdown complete.")
}
