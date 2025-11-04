package postgres

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
		SELECT
			ex.id AS id,
			ex.name AS name,
			ex.amount_cents amount,
			ex.refund_amount_cents AS refund_amount,
			ex.description AS description,
			ex.group_id AS group_id,
			cat.id AS category_id,
			ex.payer_id AS payer_id,
			ex.receiver_id AS receiver_id,
			ex.split_ratio AS split_ratio,
			ex.split_type AS split_type,
			ex.created_at AS created_at,
			ex.updated_at AS updated_at,
			ex.deleted_at AS deleted_at
		FROM expenses_latest ex INNER JOIN categories cat ON ex.category_id = cat.id
		WHERE ex.group_id = $1
		AND (ex.created_at < $2 OR (ex.created_at = $2 AND ex.id < $3))
		AND ex.document_search @@ websearch_to_tsquery('portuguese', $5)
		AND ex.deleted_at IS NULL
		ORDER BY ex.created_at DESC, ex.id DESC
		LIMIT $4
		`
	expensesQuery = `
		SELECT
			ex.id AS id,
			ex.name AS name,
			ex.amount_cents amount,
			ex.refund_amount_cents AS refund_amount,
			ex.description AS description,
			ex.group_id AS group_id,
			cat.id AS category_id,
			ex.payer_id AS payer_id,
			ex.receiver_id AS receiver_id,
			ex.split_ratio AS split_ratio,
			ex.split_type AS split_type,
			ex.created_at AS created_at,
			ex.updated_at AS updated_at,
			ex.deleted_at AS deleted_at
		FROM expenses_latest ex INNER JOIN categories cat ON ex.category_id = cat.id
		WHERE ex.group_id = $1
		AND (ex.created_at < $2 OR (ex.created_at = $2 AND ex.id < $3))
		AND ex.deleted_at IS NULL
		ORDER BY ex.created_at DESC, ex.id DESC
		LIMIT $4
  `
)

func NewGetExpenses(db *db.Client) GetExpenses {
	dbClient := db.Conn()
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
