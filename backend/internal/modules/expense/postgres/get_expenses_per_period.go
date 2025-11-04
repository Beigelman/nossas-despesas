package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

type (
	ExpensesPerPeriod struct {
		Date   string `db:"date" json:"date"`
		Amount int    `db:"amount" json:"amount"`
		Count  int    `db:"quantity" json:"quantity"`
	}

	GetExpensesPerPeriodInput struct {
		GroupID   int       `json:"group_id"`
		Aggregate string    `json:"aggregate"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	GetExpensesPerPeriod func(ctx context.Context, params GetExpensesPerPeriodInput) ([]ExpensesPerPeriod, error)
)

func NewGetExpensesPerPeriod(db *db.Client) GetExpensesPerPeriod {
	dbClient := db.Conn()
	return func(ctx context.Context, params GetExpensesPerPeriodInput) ([]ExpensesPerPeriod, error) {
		var expensesPerPeriod []ExpensesPerPeriod

		trunc, format := "day", "YYYY-MM-DD"
		switch params.Aggregate {
		case "month":
			trunc = "month"
			format = "YYYY-MM"
		case "day":
			trunc = "day"
			format = "YYYY-MM-DD"
		}

		query := fmt.Sprintf(`
			SELECT 
				to_char(date_trunc('%s', ex.created_at), '%s') AS date, 
				SUM(ex.amount_cents) AS amount, 
				COUNT(1) AS quantity 
			FROM expenses_latest ex
			WHERE ex.group_id = $1
			AND ex.created_at >= $2
			AND ex.created_at <= $3
			AND ex.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1;
		`, trunc, format)

		if err := dbClient.SelectContext(ctx, &expensesPerPeriod, query, params.GroupID, params.StartDate, params.EndDate); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return expensesPerPeriod, nil
	}
}
