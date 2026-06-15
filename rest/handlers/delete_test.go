package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type testDeleteDB struct {
	deleteErr error
}

func (db *testDeleteDB) DeleteTicket(ctx context.Context, id int64) error {
	return db.deleteErr
}

func TestDeleteTicketHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testDeleteDB{}
	router.DELETE("/tickets/:id", DeleteTicketHandler(db))

	req, _ := http.NewRequest("DELETE", "/tickets/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteTicketHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testDeleteDB{}
	router.DELETE("/tickets/:id", DeleteTicketHandler(db))

	req, _ := http.NewRequest("DELETE", "/tickets/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestDeleteTicketHandler_NotFound checks that deleting a non-existent ticket
// returns 404. The handler maps sql.ErrNoRows (returned by the DB when no row
// was deleted) to StatusNotFound.
func TestDeleteTicketHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testDeleteDB{deleteErr: sql.ErrNoRows}
	router.DELETE("/tickets/:id", DeleteTicketHandler(db))

	req, _ := http.NewRequest("DELETE", "/tickets/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteTicketHandler_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &testDeleteDB{deleteErr: context.DeadlineExceeded}
	router.DELETE("/tickets/:id", DeleteTicketHandler(db))

	req, _ := http.NewRequest("DELETE", "/tickets/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
