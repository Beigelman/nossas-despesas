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

func NewGetExpensesPerPeriod(db db.Database) GetExpensesPerPeriod {
	dbClient := db.Client()
	return func(ctx context.Context, params GetExpensesPerPeriodInput) ([]ExpensesPerPeriod, error) {
		var expensesPerPeriod []ExpensesPerPeriod

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
		        distinct on (ex.id) ex.id as id,
		        ex.amount_cents amount,
				ex.category_id  as category_id,
		        ex.created_at as created_at,
		        ex.deleted_at as deleted_at
		    from expenses ex
		    	where ex.group_id = $1
		    	and ex.created_at >= $2
		    	and ex.created_at <= $3
		    	order by ex.id desc, ex.version desc
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

		if err := dbClient.SelectContext(ctx, &expensesPerPeriod, query, params.GroupID, params.StartDate, params.EndDate); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return expensesPerPeriod, nil
	}
}
