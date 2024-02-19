package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
)

type User struct {
	ID             int     `db:"id" json:"id"`
	Name           string  `db:"name" json:"name"`
	Email          string  `db:"email" json:"email"`
	GroupID        *int    `db:"group_id" json:"group_id,omitempty"`
	ProfilePicture *string `db:"profile_picture" json:"profile_picture,omitempty"`
	CreatedAt      string  `db:"created_at" json:"created_at"`
	UpdatedAt      string  `db:"updated_at" json:"updated_at"`
}

type GetUserByID func(ctx context.Context, userID int) (*User, error)

func NewGetUserByID(db db.Database) GetUserByID {
	dbClient := db.Client()
	return func(ctx context.Context, userID int) (*User, error) {
		var user User
		if err := dbClient.GetContext(ctx, &user, `
			select id, name, email, profile_picture, group_id, created_at, updated_at
			from users
			where id = $1	
		`, userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, except.NotFoundError("user not found")
			}
			return nil, fmt.Errorf("db.GetContext: %w", err)
		}

		return &user, nil
	}
}
