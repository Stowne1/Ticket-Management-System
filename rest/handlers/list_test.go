package handlers

import (
	"Ticket-Management-System-1/postgres"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type testListDB struct {
	tickets  []postgres.Ticket
	listErr  error
	total    int
	countErr error
}

func (db *testListDB) ListTickets(ctx context.Context, limit, offset int) ([]postgres.Ticket, error) {
	if db.listErr != nil {
		return nil, db.listErr
	}
	return db.tickets, nil
}

// CountTickets satisfies the TicketLister interface, which now requires both
// ListTickets and CountTickets so the handler can return pagination metadata.
func (db *testListDB) CountTickets(ctx context.Context) (int, error) {
	if db.countErr != nil {
		return 0, db.countErr
	}
	return db.total, nil
}

func TestListTicketsHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testListDB{
		tickets: []postgres.Ticket{
			{ID: 1, Title: "First", Description: "First desc", Status: "open"},
			{ID: 2, Title: "Second", Description: "Second desc", Status: "closed"},
		},
		total: 2,
	}
	router.GET("/tickets", ListTicketsHandler(db))

	req, _ := http.NewRequest("GET", "/tickets", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	// The response is now a wrapper object: {tickets, total, page, limit}
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	tickets := body["tickets"].([]interface{})
	if len(tickets) != 2 {
		t.Errorf("Expected 2 tickets, got %d", len(tickets))
	}
	if int(body["total"].(float64)) != 2 {
		t.Errorf("Expected total 2, got %v", body["total"])
	}
}

func TestListTicketsHandler_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testListDB{tickets: []postgres.Ticket{}, total: 0}
	router.GET("/tickets", ListTicketsHandler(db))

	req, _ := http.NewRequest("GET", "/tickets", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	tickets := body["tickets"].([]interface{})
	if len(tickets) != 0 {
		t.Errorf("Expected 0 tickets, got %d", len(tickets))
	}
}

func TestListTicketsHandler_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testListDB{listErr: errors.New("db error")}
	router.GET("/tickets", ListTicketsHandler(db))

	req, _ := http.NewRequest("GET", "/tickets", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
