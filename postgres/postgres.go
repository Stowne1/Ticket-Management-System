package postgres

import (
	"database/sql"

	"context"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	_ "github.com/lib/pq"
)

type Ticket struct {
	ID          int64  `bun:"id,pk,autoincrement" json:"id"`
	Title       string `bun:",notnull" json:"title"`
	Description string `bun:",notnull" json:"description"`
	Status      string `bun:",notnull" json:"status"`
}

//https://bun.uptrace.dev/guide/models.html#mapping-tables-to-structs
//define a struct for your postgres model

type DB struct {
	Conn *bun.DB //switch to db *bun.DB
}

// Exported function
func NewDB(connStr string) (*DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return &DB{Conn: db}, nil
}

// https://bun.uptrace.dev/guide/query-insert.html
func (db *DB) InsertTicket(ctx context.Context, ticket *Ticket) error {
	_, err := db.Conn.NewInsert().Model(ticket).Exec(ctx)
	return err
}

func (db *DB) DeleteTicket(ctx context.Context, id int64) error {
	_, err := db.Conn.NewDelete().Model((*Ticket)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// https://bun.uptrace.dev/guide/query-select.html#example
func (db *DB) Retrieve(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Conn.Query(query, args...)
}

// https://bun.uptrace.dev/guide/query-update.html#example
func (db *DB) Update(query string, args ...interface{}) error {
	_, err := db.Conn.Exec(query, args...)
	return err
}

// https://bun.uptrace.dev/guide/query-delete.html
func (db *DB) Delete(query string, args ...interface{}) error {
	_, err := db.Conn.Exec(query, args...)
	return err
}

func (db *DB) UpdateTicket(ctx context.Context, ticket *Ticket) error {
	_, err := db.Conn.NewUpdate().Model(ticket).WherePK().Exec(ctx)
	return err
}
