package handlers

import (
	"Ticket-Management-System-1/postgres"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type testRegisterDB struct {
	createErr error
}

func (db *testRegisterDB) CreateUser(ctx context.Context, user *postgres.User) error {
	return db.createErr
}

func TestRegisterHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", RegisterHandler(&testRegisterDB{}))

	body := map[string]string{"email": "test@example.com", "password": "secret123"}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", RegisterHandler(&testRegisterDB{}))

	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterHandler_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", RegisterHandler(&testRegisterDB{}))

	// Omit the password field
	body := map[string]string{"email": "test@example.com"}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
