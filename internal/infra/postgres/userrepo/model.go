package userrepo

import (
	"database/sql"
	"time"
)

type UserModel struct {
	ID             int            `db:"id"`
	Name           string         `db:"name"`
	Email          string         `db:"email"`
	ProfilePicture sql.NullString `db:"profile_picture"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
	DeletedAt      sql.NullTime   `db:"deleted_at"`
	Version        int            `db:"version"`
}
