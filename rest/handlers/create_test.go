package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// MockDB implements the database interface for testing
type MockDB struct {
	shouldError bool
}

func (m *MockDB) Insert(query string, args ...interface{}) error {
	if m.shouldError {
		return &mockError{message: "database error"}
	}
	return nil
}

func (m *MockDB) Retrieve(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (m *MockDB) Update(query string, args ...interface{}) error {
	return nil
}

func (m *MockDB) Delete(query string, args ...interface{}) error {
	return nil
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

func TestCreateTicketHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: false}

	router.POST("/tickets", CreateTicketHandler(mockDB))

	ticket := Ticket{
		Title:       "Test Ticket",
		Description: "This is a test ticket",
		Status:      "open",
	}

	jsonData, _ := json.Marshal(ticket)
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["message"] != "Ticket created successfully" {
		t.Errorf("Expected success message, got %s", response["message"])
	}
}

func TestCreateTicketHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: false}

	router.POST("/tickets", CreateTicketHandler(mockDB))

	// Send invalid JSON
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBufferString(`{"invalid": json}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateTicketHandler_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: true}

	router.POST("/tickets", CreateTicketHandler(mockDB))

	ticket := Ticket{
		Title:       "Test Ticket",
		Description: "This is a test ticket",
		Status:      "open",
	}

	jsonData, _ := json.Marshal(ticket)
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
