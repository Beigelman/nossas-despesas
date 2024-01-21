package query

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/expenserepo"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"time"
)

type (
	ExpenseDetails struct {
		ID          int                    `db:"id" json:"id"`
		Name        string                 `db:"name" json:"name"`
		Amount      float32                `db:"amount" json:"amount"`
		Description string                 `db:"description" json:"description"`
		CategoryID  int                    `db:"category_id" json:"category_id"`
		PayerID     int                    `db:"payer_id" json:"payer_id"`
		ReceiverID  int                    `db:"receiver_id" json:"receiver_id"`
		SplitRatio  expenserepo.SplitRatio `db:"split_ratio" json:"split_ratio"`
		CreatedAt   time.Time              `db:"created_at" json:"created_at"`
	}

	GetGroupExpenses func(ctx context.Context, input GetGroupExpensesInput) ([]ExpenseDetails, error)

	GetGroupExpensesInput struct {
		GroupID         int
		LastExpenseDate time.Time
		Limit           int
	}
)

func NewGetGroupExpenses(db db.Database) GetGroupExpenses {
	dbClient := db.Client()
	return func(ctx context.Context, input GetGroupExpensesInput) ([]ExpenseDetails, error) {
		var expenses []ExpenseDetails
		if err := dbClient.SelectContext(ctx, &expenses, `
			with base as (
				select
    				distinct on (ex.id) ex.id as id,
    				ex.name as name,
    				ex.amount_cents amount,
    				ex.description as description,
    				cat.id as category_id,
    				payer.id as payer_id,
    				receiver.id as receiver_id,
    				ex.split_ratio as split_ratio,
					ex.created_at as created_at,
					ex.deleted_at as deleted_at
				from expenses ex
         		inner join categories cat on ex.category_id = cat.id
         		inner join users payer on ex.payer_id = payer.id
         		inner join users receiver on ex.receiver_id = receiver.id
				where ex.group_id = $1
				and ex.created_at < $2
				order by ex.id desc, ex.version desc
			)
			select id, name, amount, description, category_id, payer_id, receiver_id, split_ratio, created_at from base b
			where b.deleted_at is null
			order by b.created_at desc
			limit $3
		`, input.GroupID, input.LastExpenseDate, input.Limit); err != nil {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return expenses, nil
	}
}
