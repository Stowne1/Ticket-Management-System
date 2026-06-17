package router

import (
	fe "Ticket-Management-System-1/frontend"
	"Ticket-Management-System-1/postgres"
	"Ticket-Management-System-1/rest/handlers"
	"Ticket-Management-System-1/rest/middleware"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Setup initialises the Gin router and registers all routes.
//
// API routes (all prefixed with /api):
//
//	GET  /api/health
//	POST /api/register
//	POST /api/login
//	POST /api/tickets          (protected)
//	GET  /api/tickets          (protected)
//	GET  /api/tickets/:id      (protected)
//	PUT  /api/tickets/:id      (protected)
//	DELETE /api/tickets/:id    (protected)
//
// Static file routes:
//
//	/assets/*  — JS/CSS bundles produced by Vite
//	/*         — serves index.html so React Router handles client-side navigation
func Setup(db *postgres.DB, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// All API routes live under /api so they don't clash with React Router paths
	// like /tickets or /login which the frontend also uses.
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		api.POST("/register", handlers.RegisterHandler(db))
		api.POST("/login", handlers.LoginHandler(db, jwtSecret))

		// Protected ticket routes — rejected with 401 without a valid JWT.
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			protected.POST("/tickets", handlers.CreateTicketHandler(db))
			protected.GET("/tickets", handlers.ListTicketsHandler(db))
			protected.GET("/tickets/:id", handlers.GetTicketHandler(db))
			protected.PUT("/tickets/:id", handlers.UpdateTicketHandler(db))
			protected.DELETE("/tickets/:id", handlers.DeleteTicketHandler(db))
		}
	}

	// Serve the built React app from the embedded frontend/dist directory.
	// fs.Sub strips the "dist" prefix so paths like "dist/index.html" become "index.html".
	distFS, err := fs.Sub(fe.FS, "dist")
	if err != nil {
		panic("failed to create sub FS from embedded frontend: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(distFS))

	// /assets/* holds the Vite-generated JS and CSS bundles.
	r.GET("/assets/*filepath", func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	// serveIndex reads index.html directly from the embedded FS and writes it
	// to the response. We avoid c.FileFromFS here because http.ServeFile
	// redirects requests for "index.html" back to "/", which creates an
	// infinite redirect loop.
	serveIndex := func(c *gin.Context) {
		data, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "frontend not built — run: cd frontend && npm run build"})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	}

	r.GET("/", serveIndex)

	// Catch-all: any path that didn't match an API or asset route gets index.html.
	// This lets React Router handle client-side navigation (/login, /register, /tickets).
	// Unknown /api/* paths still get a JSON 404.
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		serveIndex(c)
	})

	return r
}
