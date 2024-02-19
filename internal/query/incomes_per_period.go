package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"time"
)

type (
	IncomesPerPeriod struct {
		Date   string `db:"date" json:"date"`
		Amount int    `db:"amount" json:"amount"`
	}

	GetIncomesPerPeriodInput struct {
		GroupID   int       `json:"group_id"`
		Aggregate string    `json:"aggregate"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	GetIncomesPerPeriod func(ctx context.Context, params GetIncomesPerPeriodInput) ([]IncomesPerPeriod, error)
)

func NewGetIncomesPerPeriod(db db.Database) GetIncomesPerPeriod {
	dbClient := db.Client()
	return func(ctx context.Context, params GetIncomesPerPeriodInput) ([]IncomesPerPeriod, error) {
		var expensesPerMonth []IncomesPerPeriod

		trunc, format := "day", "YYYY-MM-DD"
		if params.Aggregate == "month" {
			trunc = "month"
			format = "YYYY-MM"
		} else if params.Aggregate == "day" {
			trunc = "day"
			format = "YYYY-MM-DD"
		}

		query := fmt.Sprintf(`
			with base as (
			    select
			        amount_cents amount,
			        created_at as created_at,
			        deleted_at as deleted_at
			    from incomes
			    	where user_id in (select id from users where group_id = $1)
			      	and created_at >= $2
			      	and created_at <= $3
			)
			select to_char(date_trunc('%s', b.created_at), '%s') as date, sum(amount) as amount from base b
			where b.deleted_at is null
			group by 1
			order by 1;
		`, trunc, format)
		if err := dbClient.SelectContext(ctx, &expensesPerMonth, query, params.GroupID, params.StartDate, params.EndDate); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return expensesPerMonth, nil
	}
}
