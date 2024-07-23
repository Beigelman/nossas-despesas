package query

import (
	"context"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

type (
	GetExpenses func(ctx context.Context, input GetExpensesInput) ([]ExpenseDetails, error)

	GetExpensesInput struct {
		GroupID         int
		LastExpenseDate time.Time
		LastExpenseID   int
		Limit           int
		Search          string
	}
)

var (
	expensesQueryWithSearch = `
      with base as (
				select
    				distinct on (ex.id) ex.id as id,
    				ex.name as name,
    				ex.amount_cents amount,
    				ex.refund_amount_cents as refund_amount,
    				ex.description as description,
    				ex.group_id as group_id,
    				cat.id as category_id,
					  cat.icon as category_icon,
    				ex.payer_id as payer_id,
    				ex.receiver_id as receiver_id,
    				ex.split_ratio as split_ratio,
            ex.split_type as split_type,
					  ex.created_at as created_at,
					  ex.updated_at as updated_at,
					  ex.deleted_at as deleted_at,
            ts_rank(ex.document_search, websearch_to_tsquery('portuguese', $5)) as rank
				from expenses ex inner join categories cat on ex.category_id = cat.id
				where ex.group_id = $1
				and (ex.created_at < $2 or (ex.created_at = $2 and ex.id < $3))
				and ex.document_search @@ websearch_to_tsquery('portuguese', $5)
				order by ex.id desc, ex.version desc
			)
			select id, name, amount, refund_amount, description, category_id, payer_id, receiver_id, group_id, split_ratio, split_type, created_at, updated_at, deleted_at from base b
			where b.deleted_at is null
			order by b.created_at desc, b.id desc, b.rank desc
			limit $4
		`
	expensesQuery = `
    with base as (
				select
    				distinct on (ex.id) ex.id as id,
    				ex.name as name,
    				ex.amount_cents amount,
    				ex.refund_amount_cents as refund_amount,
    				ex.description as description,
    				ex.group_id as group_id,
    				cat.id as category_id,
					  cat.icon as category_icon,
    				ex.payer_id as payer_id,
    				ex.receiver_id as receiver_id,
    				ex.split_ratio as split_ratio,
            ex.split_type as split_type,
					  ex.created_at as created_at,
					  ex.updated_at as updated_at,
					  ex.deleted_at as deleted_at
				from expenses ex inner join categories cat on ex.category_id = cat.id
				where ex.group_id = $1
				and (ex.created_at < $2 or (ex.created_at = $2 and ex.id < $3))
				order by ex.id desc, ex.version desc
			)
			select id, name, amount, refund_amount, description, category_id, payer_id, receiver_id, group_id, split_ratio, split_type, created_at, updated_at, deleted_at from base b
			where b.deleted_at is null
			order by b.created_at desc, b.id desc
			limit $4
  `
)

func NewGetExpenses(db db.Database) GetExpenses {
	dbClient := db.Client()
	return func(ctx context.Context, input GetExpensesInput) ([]ExpenseDetails, error) {
		var expenses []ExpenseDetails

		query := expensesQuery
		args := []any{input.GroupID, input.LastExpenseDate, input.LastExpenseID, input.Limit}

		if input.Search != "" {
			query = expensesQueryWithSearch
			args = append(args, input.Search)
		}
		if err := dbClient.SelectContext(ctx, &expenses, query, args...); err != nil {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return expenses, nil
	}
}
