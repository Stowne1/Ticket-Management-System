package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TicketDeleter defines the interface for deleting a ticket from the database.
type TicketDeleter interface {
	DeleteTicket(ctx context.Context, id int64) error
}

// DeleteTicketHandler returns a Gin handler for deleting a ticket by its ID.
// It expects the ticket ID as a URL parameter.
// On success, it deletes the ticket from the database and returns a 200 status.
func DeleteTicketHandler(db TicketDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the ticket ID from the URL parameter
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			// If the ID is not a valid integer, return 400
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}
		// Delete the ticket using the provided database interface
		if err := db.DeleteTicket(c.Request.Context(), id); err != nil {
			// If deletion fails, return 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
			return
		}
		// On success, return a 200 status and success message
		c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
	}
}
