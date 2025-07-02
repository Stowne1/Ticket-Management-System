package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteTicketHandler handles DELETE /tickets/:id
func DeleteTicketHandler(db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}

		query := "DELETE FROM tickets WHERE id = $1"
		err = db.Delete(query, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
	}
}
