package main

import (
	"fmt"
	"os"

	"Ticket-Management-System-1/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	// Step 1: Load DATABASE_URL from env
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		fmt.Println("DATABASE_URL environment variable is not set")
		return
	}

	// Step 2: Connect to the database
	db, err := postgres.NewDB(connStr)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	// Step 3: Set up Gin
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, world!"})
	})

	router.GET("/example", func(c *gin.Context) {
		rows, err := db.Retrieve("SELECT id, name FROM your_table LIMIT 1")
		if err != nil {
			c.JSON(500, gin.H{"error": "DB error"})
			return
		}
		defer rows.Close()

		var id int
		var name string
		if rows.Next() {
			rows.Scan(&id, &name)
			c.JSON(200, gin.H{"id": id, "name": name})
		} else {
			c.JSON(404, gin.H{"error": "No data found"})
		}
	})

	// Step 4: Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
