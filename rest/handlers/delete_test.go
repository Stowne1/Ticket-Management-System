package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDeleteTicketHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: false}

	router.DELETE("/tickets/:id", DeleteTicketHandler(mockDB))

	req, _ := http.NewRequest("DELETE", "/tickets/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["message"] != "Ticket deleted successfully" {
		t.Errorf("Expected success message, got %s", response["message"])
	}
}

func TestDeleteTicketHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: false}

	router.DELETE("/tickets/:id", DeleteTicketHandler(mockDB))

	// Test with non-numeric ID
	req, _ := http.NewRequest("DELETE", "/tickets/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteTicketHandler_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: true}

	router.DELETE("/tickets/:id", DeleteTicketHandler(mockDB))

	req, _ := http.NewRequest("DELETE", "/tickets/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
