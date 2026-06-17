package handlers

import (
	"Ticket-Management-System-1/auth"
	"Ticket-Management-System-1/postgres"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserFetcher defines the DB operation needed by the login endpoint.
type UserFetcher interface {
	GetUserByEmail(ctx context.Context, email string) (*postgres.User, error)
}

// LoginHandler handles POST /login.
// On success it returns a signed JWT the client must include on every
// subsequent request as: Authorization: Bearer <token>
func LoginHandler(db UserFetcher, jwtSecret string) gin.HandlerFunc {
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

		user, err := db.GetUserByEmail(c.Request.Context(), body.Email)
		if errors.Is(err, sql.ErrNoRows) {
			// Return 401 (not 404) so we don't reveal whether this email exists.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// bcrypt.CompareHashAndPassword returns non-nil if the password is wrong.
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
			// Same 401 message as "email not found" — prevents user enumeration.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Build the JWT. The token carries the user's ID and email so downstream
		// handlers can identify the caller without hitting the DB again.
		claims := auth.Claims{
			UserID: user.ID,
			Email:  user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": signed})
	}
}
