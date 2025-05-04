package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

type (
	ExpenseDetails struct {
		ID           int                 `db:"id" json:"id"`
		Name         string              `db:"name" json:"name"`
		Amount       float32             `db:"amount" json:"amount"`
		RefundAmount *float32            `db:"refund_amount" json:"refund_amount"`
		Description  string              `db:"description" json:"description"`
		CategoryID   int                 `db:"category_id" json:"category_id"`
		PayerID      int                 `db:"payer_id" json:"payer_id"`
		ReceiverID   int                 `db:"receiver_id" json:"receiver_id"`
		GroupID      int                 `db:"group_id" json:"group_id"`
		SplitRatio   postgres.SplitRatio `db:"split_ratio" json:"split_ratio"`
		SplitType    string              `db:"split_type" json:"split_type"`
		CreatedAt    time.Time           `db:"created_at" json:"created_at"`
		UpdatedAt    time.Time           `db:"updated_at" json:"updated_at"`
		DeletedAt    *time.Time          `db:"deleted_at" json:"deleted_at"`
	}

	GetExpenseDetails func(ctx context.Context, expenseID int) ([]ExpenseDetails, error)
)

func NewGetExpenseDetails(db db.Database) GetExpenseDetails {
	dbClient := db.Client()
	return func(ctx context.Context, expenseID int) ([]ExpenseDetails, error) {
		var expenseDetails []ExpenseDetails
		if err := dbClient.SelectContext(ctx, &expenseDetails, `
			    select
    				id,
    				name,
    				amount_cents as amount,
    				refund_amount_cents as refund_amount,
    				description,
    				payer_id,
    				group_id,
    				receiver_id,
  					category_id,
    				split_ratio,
            split_type,
	  				created_at,
		  			updated_at,
			  		deleted_at
			    from expenses
				  where id = $1
			  	order by version 
		`, expenseID); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return expenseDetails, nil
	}
}
