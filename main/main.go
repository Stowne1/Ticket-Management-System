package main

import (
	"fmt"
	"os"
	"os/user"

	"Ticket-Management-System-1/postgres"
	"Ticket-Management-System-1/rest/router"
)

func main() {
	// Get current user for database connection
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	// Use current username for database connection
	connStr := fmt.Sprintf("postgres://%s@localhost:5432/ticket_system?sslmode=disable", currentUser.Username)

	// Step 2: Connect to the database
	db, err := postgres.NewDB(connStr)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	engine := router.Setup(db)
	// Step 4: Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	engine.Run(":" + port)
}
