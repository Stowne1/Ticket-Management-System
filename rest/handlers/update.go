package handlers

import (
	"Ticket-Management-System-1/postgres"
	"context"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

// TicketUpdater defines the interface for updating a ticket in the database.
type TicketUpdater interface {
	UpdateTicket(ctx context.Context, ticket *postgres.Ticket) error
}

// UpdateTicketHandler returns a Gin handler for updating a ticket by its ID.
// It expects the ticket ID as a URL parameter and a JSON body with updated fields.
// On success, it updates the ticket in the database and returns a 200 status.
func UpdateTicketHandler(db TicketUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the ticket ID from the URL parameter
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			// If the ID is not a valid integer, return 400
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}
		var ticket postgres.Ticket
		// Bind the incoming JSON to the ticket struct
		if err := c.ShouldBindJSON(&ticket); err != nil {
			// If the JSON is invalid, return 400
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		// Set the ticket ID from the URL parameter
		ticket.ID = id
		// Update the ticket using the provided database interface
		if err := db.UpdateTicket(c.Request.Context(), &ticket); err != nil {
			// If update fails, return 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
			return
		}
		// On success, return a 200 status and success message
		c.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully"})
	}
}
