package router

import (
	"Ticket-Management-System-1/postgres"
	"Ticket-Management-System-1/rest/handlers"
	"Ticket-Management-System-1/rest/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup initialises the Gin router and registers all routes.
//
// Public routes (no token required):
//
//	POST /register  — create an account
//	POST /login     — get a JWT
//	GET  /health    — liveness probe
//
// Protected routes (Authorization: Bearer <token> required):
//
//	POST   /tickets
//	GET    /tickets
//	GET    /tickets/:id
//	PUT    /tickets/:id
//	DELETE /tickets/:id
func Setup(db *postgres.DB, jwtSecret string) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes are public — the client doesn't have a token yet.
	r.POST("/register", handlers.RegisterHandler(db))
	r.POST("/login", handlers.LoginHandler(db, jwtSecret))

	// All ticket routes sit behind the auth middleware.
	// Any request without a valid JWT is rejected with 401 before reaching
	// the handler.
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.POST("/tickets", handlers.CreateTicketHandler(db))
		protected.GET("/tickets", handlers.ListTicketsHandler(db))
		protected.GET("/tickets/:id", handlers.GetTicketHandler(db))
		protected.PUT("/tickets/:id", handlers.UpdateTicketHandler(db))
		protected.DELETE("/tickets/:id", handlers.DeleteTicketHandler(db))
	}

	return r
}
