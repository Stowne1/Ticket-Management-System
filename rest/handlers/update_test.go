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

type testUpdateDB struct {
	updateErr error
}

func (db *testUpdateDB) UpdateTicket(ctx context.Context, ticket *postgres.Ticket) error {
	return db.updateErr
}

func TestUpdateTicketHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testUpdateDB{}
	router.PUT("/tickets/:id", UpdateTicketHandler(db))

	ticket := postgres.Ticket{
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
}

func TestUpdateTicketHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testUpdateDB{}
	router.PUT("/tickets/:id", UpdateTicketHandler(db))

	ticket := postgres.Ticket{
		Title:       "Updated Ticket",
		Description: "This is an updated ticket",
		Status:      "closed",
	}
	jsonData, _ := json.Marshal(ticket)
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
	db := &testUpdateDB{}
	router.PUT("/tickets/:id", UpdateTicketHandler(db))

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
	db := &testUpdateDB{updateErr: context.DeadlineExceeded}
	router.PUT("/tickets/:id", UpdateTicketHandler(db))

	ticket := postgres.Ticket{
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
