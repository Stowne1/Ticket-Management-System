package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Ticket struct should match the one in create.go
// If you move it to a shared file, import from there

// GetTicketHandler handles GET /tickets/:id
func GetTicketHandler(db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}

		query := "SELECT title, description, status FROM tickets WHERE id = $1"
		rows, err := db.Retrieve(query, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer rows.Close()

		if rows.Next() {
			var ticket Ticket
			if err := rows.Scan(&ticket.Title, &ticket.Description, &ticket.Status); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan ticket"})
				return
			}
			c.JSON(http.StatusOK, ticket)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
	}
}
