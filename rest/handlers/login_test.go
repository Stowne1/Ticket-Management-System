package handlers

import (
	"Ticket-Management-System-1/postgres"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type testLoginDB struct {
	user    *postgres.User
	findErr error
}

func (db *testLoginDB) GetUserByEmail(ctx context.Context, email string) (*postgres.User, error) {
	if db.findErr != nil {
		return nil, db.findErr
	}
	return db.user, nil
}

const testJWTSecret = "test-secret-key"

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Pre-hash a known password so we can test the full bcrypt compare path.
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), 4) // cost 4 for fast tests
	db := &testLoginDB{user: &postgres.User{ID: 1, Email: "test@example.com", PasswordHash: string(hash)}}
	router.POST("/login", LoginHandler(db, testJWTSecret))

	body := map[string]string{"email": "test@example.com", "password": "correctpassword"}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == nil || resp["token"] == "" {
		t.Errorf("Expected a token in the response, got: %v", resp)
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), 4)
	db := &testLoginDB{user: &postgres.User{ID: 1, Email: "test@example.com", PasswordHash: string(hash)}}
	router.POST("/login", LoginHandler(db, testJWTSecret))

	body := map[string]string{"email": "test@example.com", "password": "wrongpassword"}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	// sql.ErrNoRows is what the DB layer returns when no user matches the email.
	db := &testLoginDB{findErr: sql.ErrNoRows}
	router.POST("/login", LoginHandler(db, testJWTSecret))

	body := map[string]string{"email": "nobody@example.com", "password": "password"}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", LoginHandler(&testLoginDB{}, testJWTSecret))

	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(`{bad}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
