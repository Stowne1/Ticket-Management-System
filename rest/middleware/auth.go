package middleware

import (
	"Ticket-Management-System-1/auth"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware returns a Gin middleware that enforces JWT authentication.
// It must be applied to any route group that should be protected.
//
// Flow:
//  1. Extract the token from the Authorization header.
//  2. Parse and validate the token (signature + expiry).
//  3. Store the caller's userID in the Gin context for downstream handlers.
//  4. Call c.Abort() on any failure so the request stops here.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// The header must be exactly "Bearer <token>".
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be: Bearer <token>"})
			c.Abort()
			return
		}

		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			// Explicitly check the signing method. Without this check an attacker
			// could craft a token signed with "none" (no signature) and it would pass.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Make the caller's userID available to every handler in this request.
		// Handlers read it with: userID := c.GetInt64("userID")
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
