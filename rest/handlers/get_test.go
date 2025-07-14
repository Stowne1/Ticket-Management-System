package handlers

import (
	"Ticket-Management-System-1/postgres"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type testGetDB struct {
	findErr error
	found   bool
	result  *postgres.Ticket
}

func (db *testGetDB) GetTicketByID(ctx context.Context, id int64) (*postgres.Ticket, error) {
	if db.findErr != nil {
		return nil, db.findErr
	}
	if !db.found {
		return nil, nil
	}
	return db.result, nil
}

func TestGetTicketHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testGetDB{}
	router.GET("/tickets/:id", GetTicketHandler(db))

	req, _ := http.NewRequest("GET", "/tickets/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetTicketHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testGetDB{found: false}
	router.GET("/tickets/:id", GetTicketHandler(db))

	req, _ := http.NewRequest("GET", "/tickets/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetTicketHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testGetDB{found: true, result: &postgres.Ticket{ID: 1, Title: "Test", Description: "Desc", Status: "open"}}
	router.GET("/tickets/:id", GetTicketHandler(db))

	req, _ := http.NewRequest("GET", "/tickets/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
