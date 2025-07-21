package handlers

import (
	"Ticket-Management-System-1/postgres"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TicketGetter defines the interface for retrieving a ticket by ID from the database.
type TicketGetter interface {
	GetTicketByID(ctx context.Context, id int64) (*postgres.Ticket, error)
}

// GetTicketHandler returns a Gin handler for fetching a ticket by its ID.
// It expects the ticket ID as a URL parameter.
// On success, it returns the ticket as JSON. If not found, returns 404.
func GetTicketHandler(db TicketGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the ticket ID from the URL parameter
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			// If the ID is not a valid integer, return 400
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}
		// Retrieve the ticket from the database
		ticket, err := db.GetTicketByID(c.Request.Context(), id)
		if err != nil {
			// If retrieval fails, return 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if ticket == nil {
			// If no ticket is found, return 404
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}
		// Return the ticket as JSON
		c.JSON(http.StatusOK, ticket)
	}
}
