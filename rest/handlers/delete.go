package handlers

import (
	"Ticket-Management-System-1/postgres"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteTicketHandler(db *postgres.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}
		if err := db.DeleteTicket(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
	}
}
