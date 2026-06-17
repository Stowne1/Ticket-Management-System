package main

import (
	"Ticket-Management-System-1/postgres"
	"Ticket-Management-System-1/rest/router"
	"Ticket-Management-System-1/migrations"
	"github.com/uptrace/bun/migrate"
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
	defer db.Conn.Close()

	// Run migrations on startup
	ctx := context.Background()
	migs := migrate.NewMigrations()
	if err := migs.Discover(migrations.FS); err != nil {
		log.Fatalf("Failed to discover migrations: %v", err)
	}
	migrator := migrate.NewMigrator(db.Conn, migs)
	if err := migrator.Init(ctx); err != nil {
		log.Fatalf("Failed to init migrator: %v", err)
	}
	if _, err := migrator.Migrate(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// JWT_SECRET is used to sign and verify all tokens.
	// The server refuses to start without it — an empty secret would let
	// anyone forge tokens by signing with an empty string.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Set up the Gin router with all ticket handlers
	r := router.Setup(db, jwtSecret)

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
