package handlers

import (
	"Ticket-Management-System-1/postgres"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketLister interface {
	ListTickets(ctx context.Context, limit, offset int) ([]postgres.Ticket, error)
}

func ListTicketsHandler(db TicketLister) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 10
		page := 1

		if l := c.Query("limit"); l != "" {
			if val, err := strconv.Atoi(l); err == nil && val > 0 {
				limit = val
			}
		}
		if p := c.Query("page"); p != "" {
			if val, err := strconv.Atoi(p); err == nil && val > 0 {
				page = val
			}
		}
		offset := (page - 1) * limit

		tickets, err := db.ListTickets(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tickets"})
			return 
		}
		c.JSON(http.StatusOK, tickets)
	}
} 