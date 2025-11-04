package postgres

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

type (
	UserBalance struct {
		UserID  int     `db:"user_id" json:"user_id"`
		Balance float32 `db:"balance" json:"balance"`
	}

	GetGroupBalance func(ctx context.Context, groupID int) ([]UserBalance, error)
)

func NewGetGroupBalance(db *db.Client) GetGroupBalance {
	dbClient := db.Conn()
	return func(ctx context.Context, groupID int) ([]UserBalance, error) {
		var balances []UserBalance
		if err := dbClient.SelectContext(ctx, &balances, `
			WITH base AS (
			    SELECT
			        id,
			        CASE WHEN refund_amount_cents IS NULL THEN amount_cents ELSE amount_cents - refund_amount_cents END AS amount_cents,
			        group_id,
			        split_ratio,
			        payer_id,
			        receiver_id,
			        deleted_at
			    FROM expenses_latest
			    WHERE group_id = $1
			    AND deleted_at IS NULL
			), balances AS (
			    SELECT
			        user_id,
			        balance,
			        type
			    FROM (
					SELECT payer_id AS user_id, SUM((amount_cents * (split_ratio->>'receiver')::numeric / 100)) AS balance, 'payer' as type
					FROM base
					GROUP BY payer_id
	
					UNION ALL
	
					SELECT receiver_id AS user_id, SUM((amount_cents * (split_ratio->>'receiver')::numeric / 100)) AS balance, 'receiver' as type
					FROM base
					GROUP BY receiver_id
			 	) AS balances
			)
			SELECT
			    user_id,
			    SUM(CASE WHEN type = 'payer' THEN balance ELSE balance * -1 END) AS balance
			FROM balances
			GROUP BY user_id
		`, &groupID); err != nil {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		return balances, nil
	}
}
