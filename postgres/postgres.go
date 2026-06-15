package postgres

import (
	"database/sql"

	"context"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	_ "github.com/lib/pq"
)

// Ticket represents a support ticket in the system.
// The struct tags configure Bun ORM and JSON serialization.
type Ticket struct {
	ID          int64     `bun:"id,pk,autoincrement" json:"id"` // Primary key, auto-incremented
	Title       string    `bun:",notnull" json:"title"`         // Title of the ticket
	Description string    `bun:",notnull" json:"description"`   // Description of the issue
	Status      string    `bun:",notnull" json:"status"`        // Status (e.g., open, closed)
	CreatedAt   time.Time `bun:",notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:",notnull,default:current_timestamp" json:"updated_at"`
}

// DB wraps a Bun database connection and provides methods for ticket operations.
type DB struct {
	Conn *bun.DB // Bun DB connection
}

// NewDB creates a new Bun DB connection using the provided connection string.
// It sets up the underlying SQL DB and configures the Bun dialect for Postgres.
func NewDB(connStr string) (*DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
	db := bun.NewDB(sqldb, pgdialect.New())
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{Conn: db}, nil
}

// InsertTicket inserts a new ticket into the database using Bun ORM.
// Returns an error if the insert fails.
func (db *DB) InsertTicket(ctx context.Context, ticket *Ticket) error {
	_, err := db.Conn.NewInsert().Model(ticket).Exec(ctx)
	return err
}

// UpdateTicket updates an existing ticket in the database using Bun ORM.
// The ticket must have a valid ID (primary key).
func (db *DB) UpdateTicket(ctx context.Context, ticket *Ticket) error {
	// Set UpdatedAt here rather than relying on the caller or the DB default.
	// The DEFAULT NOW() in the schema only fires on INSERT, not UPDATE, so
	// without this line every update would write Go's zero time (0001-01-01).
	ticket.UpdatedAt = time.Now()
	result, err := db.Conn.NewUpdate().Model(ticket).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil

}

// DeleteTicket deletes a ticket by its ID using Bun ORM.
// Returns an error if the delete fails or the ticket does not exist.
func (db *DB) DeleteTicket(ctx context.Context, id int64) error {
	result, err := db.Conn.NewDelete().Model((*Ticket)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetTicketByID retrieves a ticket by its ID using Bun ORM.
// Returns the ticket if found, or an error if not found or on DB error.
func (db *DB) GetTicketByID(ctx context.Context, id int64) (*Ticket, error) {
	ticket := new(Ticket)
	err := db.Conn.NewSelect().
		Model(ticket).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (db *DB) ListTickets(ctx context.Context, limit, offset int) ([]Ticket, error) {
	// Initialise as an empty (non-nil) slice so JSON serialises to [] not null
	// when there are no tickets.
	tickets := make([]Ticket, 0)
	err := db.Conn.NewSelect().
		Model(&tickets).
		// ORDER BY id ensures the order is deterministic across pages.
		// Without this, Postgres can return rows in any order.
		OrderExpr("id ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

// CountTickets returns the total number of tickets in the database.
// Used by the list endpoint to tell clients how many pages exist.
func (db *DB) CountTickets(ctx context.Context) (int, error) {
	count, err := db.Conn.NewSelect().Model((*Ticket)(nil)).Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}
