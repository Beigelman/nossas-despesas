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
		Category    string                 `db:"category" json:"category"`
		Payer       string                 `db:"payer" json:"payer"`
		Receiver    string                 `db:"receiver" json:"receiver"`
		SplitRatio  expenserepo.SplitRatio `db:"split_ratio" json:"split_ratio"`
		CreatedAt   time.Time              `db:"created_at" json:"created_at"`
	}

	GetGroupExpenses func(ctx context.Context, input GetGroupExpensesInput) ([]ExpenseDetails, error)

	GetGroupExpensesInput struct {
		GroupID       int
		LastExpenseID int
		Limit         int
	}
)

func NewGetGroupExpenses(db db.Database) GetGroupExpenses {
	dbClient := db.Client()
	return func(ctx context.Context, input GetGroupExpensesInput) ([]ExpenseDetails, error) {
		var expenses []ExpenseDetails

		if err := dbClient.SelectContext(ctx, &expenses, `
			select
    			distinct on (ex.id) ex.id as id,
    			ex.name as name,
    			(ex.amount_cents::float / 100) as amount,
    			ex.description as description,
    			cat.name as category,
    			payer.name as payer,
    			receiver.name as receiver,
    			ex.split_ratio as split_ratio,
				ex.created_at as created_at
			from expenses ex
         	inner join categories cat on ex.category_id = cat.id
         	inner join users payer on ex.payer_id = payer.id
         	inner join users receiver on ex.receiver_id = receiver.id
			where ex.group_id = $1
			and ex.id > $2
			and ex.deleted_at is null
			order by ex.id, ex.version desc
			limit $3
		`, input.GroupID, input.LastExpenseID, input.Limit); err != nil {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return expenses, nil
	}
}
