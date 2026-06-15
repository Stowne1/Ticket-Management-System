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

// TestUpdateTicketHandler_NotFound checks that updating a non-existent ticket
// returns 404. The DB returns sql.ErrNoRows when WherePK matches zero rows.
func TestUpdateTicketHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testUpdateDB{updateErr: sql.ErrNoRows}
	router.PUT("/tickets/:id", UpdateTicketHandler(db))

	ticket := postgres.Ticket{
		Title:       "Ghost Ticket",
		Description: "This ticket does not exist",
		Status:      "open",
	}
	jsonData, _ := json.Marshal(ticket)
	req, _ := http.NewRequest("PUT", "/tickets/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
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
