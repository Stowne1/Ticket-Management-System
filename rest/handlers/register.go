package handlers

import (
	"Ticket-Management-System-1/postgres"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/bcrypt"
)

// UserCreator defines the DB operation needed by the register endpoint.
type UserCreator interface {
	CreateUser(ctx context.Context, user *postgres.User) error
}

// RegisterHandler handles POST /register.
// It hashes the password with bcrypt before storing it — we never write
// plaintext passwords to the database.
func RegisterHandler(db UserCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		if body.Email == "" || body.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
			return
		}

		// bcrypt cost 12 is the recommended minimum for production. It makes
		// brute-forcing a stolen hash database expensive even with modern hardware.
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := &postgres.User{
			Email:        body.Email,
			PasswordHash: string(hash),
		}

		if err := db.CreateUser(c.Request.Context(), user); err != nil {
			// Postgres error code 23505 is "unique_violation". The email column
			// has a UNIQUE constraint, so this fires when the email is taken.
			// We return 409 Conflict rather than a generic 500 so the client
			// knows to prompt the user to log in instead.
			var pgErr pgdriver.Error
			if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			"id":      user.ID,
		})
	}
}
