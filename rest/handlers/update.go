package handlers

//write a put api handeler that updates a ticket by its id

import (
	"Ticket-Management-System-1/postgres"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateTicketHandler(db *postgres.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}
		var ticket postgres.Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		ticket.ID = id
		if err := db.UpdateTicket(c.Request.Context(), &ticket); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully"})
	}
}
