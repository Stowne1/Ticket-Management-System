package main

import (
	"Ticket-Management-System-1/postgres"
	"Ticket-Management-System-1/rest/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main is the entry point for the Ticket Management System server.
// It initializes the database, sets up the router, and starts the HTTP server.
func main() {
	// Get the Postgres connection string from the environment variable
	connStr := os.Getenv("POSTGRES_DSN")
	if connStr == "" {
		log.Fatal("POSTGRES_DSN environment variable is not set")
	}

	// Initialize the Bun-backed Postgres DB
	db, err := postgres.NewDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set up the Gin router with all ticket handlers
	r := router.Setup(db)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
