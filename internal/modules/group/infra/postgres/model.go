package postgres

import (
	"database/sql"
	"time"
)

type GroupModel struct {
	ID        int          `db:"id"`
	Name      string       `db:"name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	Version   int          `db:"version"`
}

type GroupInviteModel struct {
	ID        int          `db:"id"`
	GroupID   int          `db:"group_id"`
	Email     string       `db:"email"`
	Status    string       `db:"status"`
	Token     string       `db:"token"`
	ExpiresAt time.Time    `db:"expires_at"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	Version   int          `db:"version"`
}
