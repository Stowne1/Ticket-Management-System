package handlers

//write a put api handeler that updates a ticket by its id

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateTicketHandler handles PUT /tickets/:id
func UpdateTicketHandler(db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}

		var ticket Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		query := "UPDATE tickets SET title = $1, description = $2, status = $3 WHERE id = $4"
		err = db.Update(query, ticket.Title, ticket.Description, ticket.Status, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully"})
	}
}
