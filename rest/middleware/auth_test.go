package middleware

import (
	"Ticket-Management-System-1/auth"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const secret = "test-secret"

// makeToken generates a signed JWT for use in tests.
// expiry controls whether the token is valid or already expired.
func makeToken(userID int64, expiry time.Time) string {
	claims := auth.Claims{
		UserID: userID,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func setupMiddlewareRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", AuthMiddleware(secret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"userID": c.GetInt64("userID")})
	})
	return r
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	token := makeToken(42, time.Now().Add(time.Hour))
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	setupMiddlewareRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	setupMiddlewareRouter().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer this.is.not.a.valid.token")
	w := httptest.NewRecorder()

	setupMiddlewareRouter().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// Create a token that expired one hour ago.
	token := makeToken(42, time.Now().Add(-time.Hour))
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	setupMiddlewareRouter().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_MalformedHeader(t *testing.T) {
	// Header present but missing the "Bearer" prefix.
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "just-a-token-no-bearer")
	w := httptest.NewRecorder()

	setupMiddlewareRouter().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
