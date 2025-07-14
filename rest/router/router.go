package router

import (
	"Ticket-Management-System-1/postgres"
	"Ticket-Management-System-1/rest/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the Gin router with all ticket handlers
func Setup(db *postgres.DB) *gin.Engine {
	router := gin.Default()
	router.POST("/tickets", handlers.CreateTicketHandler(db))
	router.GET("/tickets/:id", handlers.GetTicketHandler(db))
	router.PUT("/tickets/:id", handlers.UpdateTicketHandler(db))
	router.DELETE("/tickets/:id", handlers.DeleteTicketHandler(db))
	return router
}
