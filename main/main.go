package main

import (
	"Ticket-Management-System-1/postgres"
	"Ticket-Management-System-1/rest/router"
	"log"
	"os"
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

	// Start the HTTP server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
