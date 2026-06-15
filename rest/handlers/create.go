package handlers

import (
	"Ticket-Management-System-1/postgres"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// validStatuses is the complete set of values the status field may hold.
// Validating here (rather than relying on the DB) gives callers a clear error
// message instead of a silent bad write or a generic 500.
var validStatuses = map[string]bool{
	"open":        true,
	"in_progress": true,
	"closed":      true,
}

// TicketInserter defines the interface for inserting a ticket into the database.
type TicketInserter interface {
	InsertTicket(ctx context.Context, ticket *postgres.Ticket) error
}

// CreateTicketHandler returns a Gin handler for creating a new ticket.
// It expects a JSON body with title, description, and status fields.
// On success, it inserts the ticket into the database and returns a 201 status.
func CreateTicketHandler(db TicketInserter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticket postgres.Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		if ticket.Title == "" || ticket.Description == "" || ticket.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}
		// Reject any status value that isn't in our allowed set.
		if !validStatuses[ticket.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status must be one of: open, in_progress, closed"})
			return
		}
		if err := db.InsertTicket(c.Request.Context(), &ticket); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Ticket created successfully",
			"id":      ticket.ID,
		})
	}
}
