package postgres

import (
	"database/sql"
	"time"
)

type AuthModel struct {
	ID         int            `db:"id"`
	Email      string         `db:"email"`
	Password   sql.NullString `db:"password"`
	ProviderID sql.NullString `db:"provider_id"`
	Type       string         `db:"type"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
	DeletedAt  sql.NullTime   `db:"deleted_at"`
	Version    int            `db:"version"`
}
