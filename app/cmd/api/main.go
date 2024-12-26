package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/AmazingAkai/URL-Shortener/app/internal/autoload"
	"github.com/AmazingAkai/URL-Shortener/app/internal/routes"
	"github.com/AmazingAkai/URL-Shortener/app/internal/server"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func main() {
	server := server.New()
	routes.RegisterURLRoutes(server)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		err := server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	done := make(chan bool, 1)
	go gracefulShutdown(server, done)

	<-done
	log.Println("Graceful shutdown complete.")
}
