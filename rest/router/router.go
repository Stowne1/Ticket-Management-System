package router

import (
	"Ticket-Management-System-1/postgres"
	"Ticket-Management-System-1/rest/handlers"

	"github.com/gin-gonic/gin"
)

// Setup initializes the Gin router and registers all ticket-related routes.
// It takes a Bun-backed DB and passes it to each handler via the appropriate interface.
func Setup(db *postgres.DB) *gin.Engine {
	// Create a new Gin router instance
	router := gin.Default()

	// Register the POST /tickets route for creating a new ticket
	router.POST("/tickets", handlers.CreateTicketHandler(db))

	// Register the GET /tickets/:id route for fetching a ticket by ID
	router.GET("/tickets/:id", handlers.GetTicketHandler(db))

	// Register the PUT /tickets/:id route for updating a ticket by ID
	router.PUT("/tickets/:id", handlers.UpdateTicketHandler(db))

	// Register the DELETE /tickets/:id route for deleting a ticket by ID
	router.DELETE("/tickets/:id", handlers.DeleteTicketHandler(db))

	// Return the configured router
	return router
}
