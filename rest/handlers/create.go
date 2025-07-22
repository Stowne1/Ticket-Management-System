package handlers

import (
	"Ticket-Management-System-1/postgres"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		// Bind the incoming JSON to the ticket struct
		if err := c.ShouldBindJSON(&ticket); err != nil {
			// If the JSON is invalid, return a 400 error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		// Validate required fields
		if ticket.Title == "" || ticket.Description == "" || ticket.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}
		// Insert the ticket using the provided database interface
		if err := db.InsertTicket(c.Request.Context(), &ticket); err != nil {
			// If insertion fails, return a 500 error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// On success, return a 201 status and success message
		c.JSON(http.StatusCreated, gin.H{
			"message": "Ticket created successfully",
			"id":      ticket.ID,
		})
	}
}
