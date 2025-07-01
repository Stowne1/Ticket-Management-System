package postgres

import (
    "database/sql"
    _ "github.com/lib/pq"
)
type DB struct {
    Conn *sql.DB
}
// Exported function
func NewDB(connStr string) (*DB, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return &DB{Conn: db}, nil
}

func (db *DB) Insert(query string, args ...interface{}) error {
    _, err := db.Conn.Exec(query, args...)
    return err
}

func (db *DB) Retrieve(query string, args ...interface{}) (*sql.Rows, error) {
    return db.Conn.Query(query, args...)
}

func (db *DB) Update(query string, args ...interface{}) error {
    _, err := db.Conn.Exec(query, args...)
    return err
}

func (db *DB) Delete(query string, args ...interface{}) error {
    _, err := db.Conn.Exec(query, args...)
    return err
}
