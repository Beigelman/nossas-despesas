package postgres

import (
	"database/sql"
	"time"
)

type IncomeModel struct {
	ID        int          `db:"id"`
	UserID    int          `db:"user_id"`
	Amount    int          `db:"amount_cents"`
	Type      string       `db:"type"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	Version   int          `db:"version"`
}
