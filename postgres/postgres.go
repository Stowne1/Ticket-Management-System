package postgres

import (
	"database/sql"

	"context"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	_ "github.com/lib/pq"
)

// Ticket represents a support ticket in the system.
// The struct tags configure Bun ORM and JSON serialization.
type Ticket struct {
	ID          int64  `bun:"id,pk,autoincrement" json:"id"` // Primary key, auto-incremented
	Title       string `bun:",notnull" json:"title"`         // Title of the ticket
	Description string `bun:",notnull" json:"description"`   // Description of the issue
	Status      string `bun:",notnull" json:"status"`        // Status (e.g., open, closed)
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
	_, err := db.Conn.NewUpdate().Model(ticket).WherePK().Exec(ctx)
	return err
}

// DeleteTicket deletes a ticket by its ID using Bun ORM.
// Returns an error if the delete fails or the ticket does not exist.
func (db *DB) DeleteTicket(ctx context.Context, id int64) error {
	_, err := db.Conn.NewDelete().Model((*Ticket)(nil)).Where("id = ?", id).Exec(ctx)
	return err
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
