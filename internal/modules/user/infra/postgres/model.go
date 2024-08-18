package postgres

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type UserModel struct {
	ID             int            `db:"id"`
	Name           string         `db:"name"`
	Email          string         `db:"email"`
	GroupID        sql.NullInt64  `db:"group_id"`
	ProfilePicture sql.NullString `db:"profile_picture"`
	Flags          pq.StringArray `db:"flags"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
	DeletedAt      sql.NullTime   `db:"deleted_at"`
	Version        int            `db:"version"`
}
