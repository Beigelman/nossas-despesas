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
	IncomesPerPeriod struct {
		Date   string `db:"date" json:"date"`
		Amount int    `db:"amount" json:"amount"`
		Count  int    `db:"quantity" json:"quantity"`
	}

	GetIncomesPerPeriodInput struct {
		GroupID   int       `json:"group_id"`
		Aggregate string    `json:"aggregate"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	GetIncomesPerPeriod func(ctx context.Context, params GetIncomesPerPeriodInput) ([]IncomesPerPeriod, error)
)

func NewGetIncomesPerPeriod(db *db.Client) GetIncomesPerPeriod {
	dbClient := db.Conn()
	return func(ctx context.Context, params GetIncomesPerPeriodInput) ([]IncomesPerPeriod, error) {
		var expensesPerMonth []IncomesPerPeriod

		trunc, format := "day", "YYYY-MM-DD"
		switch params.Aggregate {
		case "month":
			trunc = "month"
			format = "YYYY-MM"
		case "day":
			trunc = "day"
			format = "YYYY-MM-DD"
		}

		// TODO: Acho que esse base é desnecessário, pois o incomes_latest já tem o deleted_at
		query := fmt.Sprintf(`
			WITH base AS (
			    SELECT
			        amount_cents amount,
			        created_at AS created_at,
			        deleted_at AS deleted_at
			    FROM incomes
			    WHERE user_id IN (SELECT id FROM users WHERE group_id = $1)
			    AND created_at >= $2
			    AND created_at <= $3
			)
			SELECT 
				to_char(date_trunc('%s', b.created_at), '%s') AS date, 
				SUM(amount) AS amount,
				COUNT(1) AS quantity
			FROM base b
			WHERE b.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1;
		`, trunc, format)
		if err := dbClient.SelectContext(
			ctx, &expensesPerMonth,
			query, params.GroupID,
			params.StartDate, params.EndDate,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return expensesPerMonth, nil
	}
}
