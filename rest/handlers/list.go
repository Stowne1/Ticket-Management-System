package handlers

import (
	"Ticket-Management-System-1/postgres"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TicketLister defines the database operations needed by the list endpoint.
// It combines fetching a page of tickets with counting the total so the
// client can calculate how many pages exist.
type TicketLister interface {
	ListTickets(ctx context.Context, limit, offset int) ([]postgres.Ticket, error)
	CountTickets(ctx context.Context) (int, error)
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

		// Fetch the total count so the client can calculate total pages.
		// Formula: total_pages = ceil(total / limit).
		total, err := db.CountTickets(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count tickets"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tickets": tickets,
			"total":   total,
			"page":    page,
			"limit":   limit,
		})
	}
}
