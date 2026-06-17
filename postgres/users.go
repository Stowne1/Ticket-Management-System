package postgres

import (
	"context"
	"time"
)

// User represents a registered account in the system.
type User struct {
	ID           int64     `bun:"id,pk,autoincrement" json:"id"`
	Email        string    `bun:",notnull,unique"     json:"email"`
	// json:"-" ensures the hash is never included in any API response,
	// even if a handler accidentally serialises the whole struct.
	PasswordHash string    `bun:"password_hash,notnull" json:"-"`
	CreatedAt    time.Time `bun:",notnull,default:current_timestamp" json:"created_at"`
}

// CreateUser inserts a new user row. The caller is responsible for hashing
// the password before passing the User in.
func (db *DB) CreateUser(ctx context.Context, user *User) error {
	_, err := db.Conn.NewInsert().Model(user).Exec(ctx)
	return err
}

// GetUserByEmail looks up a user by their email address.
// Returns sql.ErrNoRows if no matching row exists.
func (db *DB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := new(User)
	err := db.Conn.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
