package main

import (
	"fmt"
	"os"
	"os/user"

	"Ticket-Management-System-1/postgres"
	"Ticket-Management-System-1/rest/handlers"

	"github.com/gin-gonic/gin"
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

	// Step 3: Set up Gin
	router := gin.Default()
	router.POST("/tickets", handlers.CreateTicketHandler(db))
	router.GET("/tickets/:id", handlers.GetTicketHandler(db))
	router.PUT("/tickets/:id", handlers.UpdateTicketHandler(db))
	router.DELETE("/tickets/:id", handlers.DeleteTicketHandler(db))

	// Step 4: Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
