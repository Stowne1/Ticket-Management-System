package handlers

//write a post api handler that calls the postgres insert method

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Database interface for dependency injection
type Database interface {
	Insert(query string, args ...interface{}) error
	Retrieve(query string, args ...interface{}) (*sql.Rows, error)
	Update(query string, args ...interface{}) error
	Delete(query string, args ...interface{}) error
}

// Ticket represents a support ticket
// ID is omitted in POST, assumed auto-incremented by DB
// Adjust fields as needed
type Ticket struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// CreateTicketHandler handles POST /tickets
func CreateTicketHandler(db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticket Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if ticket.Title == "" || ticket.Description == "" || ticket.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		query := "INSERT INTO tickets (title, description, status) VALUES ($1, $2, $3)"
		err := db.Insert(query, ticket.Title, ticket.Description, ticket.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert ticket"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Ticket created successfully"})
	}
}
