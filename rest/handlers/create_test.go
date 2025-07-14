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

type testDB struct {
	insertErr error
}

func (db *testDB) InsertTicket(ctx context.Context, ticket *postgres.Ticket) error {
	return db.insertErr
}

func TestCreateTicketHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testDB{}
	router.POST("/tickets", CreateTicketHandler(db))

	ticket := postgres.Ticket{
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
}

func TestCreateTicketHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testDB{}
	router.POST("/tickets", CreateTicketHandler(db))

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
	db := &testDB{insertErr: context.DeadlineExceeded}
	router.POST("/tickets", CreateTicketHandler(db))

	ticket := postgres.Ticket{
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
