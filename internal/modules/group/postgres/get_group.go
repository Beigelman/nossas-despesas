package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	Member struct {
		ID             int       `db:"id" json:"id"`
		Name           string    `db:"name" json:"name"`
		Email          string    `db:"email" json:"email"`
		GroupID        int       `db:"group_id" json:"group_id"`
		ProfilePicture *string   `db:"profile_picture" json:"profile_picture,omitempty"`
		CreatedAt      time.Time `db:"created_at" json:"created_at"`
		UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
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

func NewGetGroup(db *db.Client) GetGroup {
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
			email,
			group_id,
			profile_picture,
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
