package postgres

import (
	"database/sql"
	"time"
)

type CategoryModel struct {
	ID              int          `db:"id"`
	Name            string       `db:"name"`
	Icon            string       `db:"icon"`
	CategoryGroupID int          `db:"category_group_id"`
	CreatedAt       time.Time    `db:"created_at"`
	UpdatedAt       time.Time    `db:"updated_at"`
	DeletedAt       sql.NullTime `db:"deleted_at"`
	Version         int          `db:"version"`
}

type CategoryGroupModel struct {
	ID        int          `db:"id"`
	Name      string       `db:"name"`
	Icon      string       `db:"icon"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	Version   int          `db:"version"`
}
