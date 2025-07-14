package handlers

import (
	"Ticket-Management-System-1/postgres"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Ticket struct should match the one in create.go
// If you move it to a shared file, import from there

// GetTicketHandler handles GET /tickets/:id
func GetTicketHandler(db *postgres.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			return
		}

		ticket := new(postgres.Ticket)
		err = db.Conn.NewSelect().
			Model(ticket).
			Where("id = ?", id).
			Scan(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ticket ID"})
			return
		}
		c.JSON(http.StatusOK, ticket)

	}
}
