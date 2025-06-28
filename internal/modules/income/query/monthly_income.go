package query

import (
	"context"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

type (
	UserIncome struct {
		ID        int       `db:"id" json:"id"`
		UserID    int       `db:"user_id" json:"user_id"`
		Type      string    `db:"type" json:"type"`
		Amount    int       `db:"amount_cents" json:"amount"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
	}

	GetMonthlyIncome func(ctx context.Context, groupID int, date time.Time) ([]UserIncome, error)
)

func NewGetMonthlyIncome(db *db.Client) GetMonthlyIncome {
	dbClient := db.Client()
	return func(ctx context.Context, groupID int, date time.Time) ([]UserIncome, error) {
		var balances []UserIncome
		if err := dbClient.SelectContext(ctx, &balances, `
			SELECT id, user_id, type, amount_cents, created_at 
			FROM incomes
			WHERE user_id IN (SELECT id FROM users WHERE group_id = $1)
			AND extract(month from created_at) = $2
			AND extract(year from created_at) = $3
			AND deleted_at is null
		`, groupID, date.Month(), date.Year()); err != nil {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return balances, nil
	}
}
