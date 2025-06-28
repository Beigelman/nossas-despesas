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
			select 
				to_char(date_trunc('%s', b.created_at), '%s') as date, 
				sum(amount) as amount,
				count(1) as quantity
			from base b
			where b.deleted_at is null
			group by 1
			order by 1;
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
