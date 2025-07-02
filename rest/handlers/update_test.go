package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUpdateTicketHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: false}

	router.PUT("/tickets/:id", UpdateTicketHandler(mockDB))

	ticket := Ticket{
		Title:       "Updated Ticket",
		Description: "This is an updated ticket",
		Status:      "closed",
	}

	jsonData, _ := json.Marshal(ticket)
	req, _ := http.NewRequest("PUT", "/tickets/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["message"] != "Ticket updated successfully" {
		t.Errorf("Expected success message, got %s", response["message"])
	}
}

func TestUpdateTicketHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: false}

	router.PUT("/tickets/:id", UpdateTicketHandler(mockDB))

	ticket := Ticket{
		Title:       "Updated Ticket",
		Description: "This is an updated ticket",
		Status:      "closed",
	}

	jsonData, _ := json.Marshal(ticket)
	// Test with non-numeric ID
	req, _ := http.NewRequest("PUT", "/tickets/abc", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateTicketHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: false}

	router.PUT("/tickets/:id", UpdateTicketHandler(mockDB))

	// Send invalid JSON
	req, _ := http.NewRequest("PUT", "/tickets/1", bytes.NewBufferString(`{"invalid": json}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateTicketHandler_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{shouldError: true}

	router.PUT("/tickets/:id", UpdateTicketHandler(mockDB))

	ticket := Ticket{
		Title:       "Updated Ticket",
		Description: "This is an updated ticket",
		Status:      "closed",
	}

	jsonData, _ := json.Marshal(ticket)
	req, _ := http.NewRequest("PUT", "/tickets/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
