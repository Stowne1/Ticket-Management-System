package handlers

//write a post api handler that calls the postgres insert method

import (
	"net/http"

	"Ticket-Management-System-1/postgres"

	"github.com/gin-gonic/gin"
)

func CreateTicketHandler(db *postgres.DB) gin.HandlerFunc {
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

		if err := db.InsertTicket(c.Request.Context(), &ticket); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert ticket"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Ticket created successfully"})
	}
}
