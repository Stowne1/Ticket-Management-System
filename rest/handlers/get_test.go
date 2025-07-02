package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetTicketHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: false}

	router.GET("/tickets/:id", GetTicketHandler(mockDB))

	// Test with non-numeric ID
	req, _ := http.NewRequest("GET", "/tickets/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetTicketHandler_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: true}

	router.GET("/tickets/:id", GetTicketHandler(mockDB))

	req, _ := http.NewRequest("GET", "/tickets/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGetTicketHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: false}

	router.GET("/tickets/:id", GetTicketHandler(mockDB))

	req, _ := http.NewRequest("GET", "/tickets/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 404 when no data is found (our mock returns empty result)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
