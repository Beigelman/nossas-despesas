package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"time"
)

type (
	Member struct {
		ID        int       `db:"id" json:"id"`
		Name      string    `db:"name" json:"name"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}

	Group struct {
		ID        int       `db:"id" json:"id"`
		Name      string    `db:"name" json:"name"`
		Members   []Member  `json:"members"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}

	GetGroup func(ctx context.Context, groupID int) (*Group, error)
)

func NewGetGroup(db db.Database) GetGroup {
	dbClient := db.Client()
	return func(ctx context.Context, groupID int) (*Group, error) {
		var group Group

		if err := dbClient.GetContext(ctx, &group, `
			select
    		id,
			name,
			created_at,
			updated_at 
			from groups
			where id = $1
		`, groupID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, except.NotFoundError("group not found")
			}
			return nil, fmt.Errorf("db.GetContext: %w", err)
		}

		var members []Member
		if err := dbClient.SelectContext(ctx, &members, `
			select
    		id,
			name,
			created_at,
			updated_at 
			from users
			where group_id = $1
		`, groupID); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.Select: %w", err)
		}

		group.Members = members

		return &group, nil
	}
}
