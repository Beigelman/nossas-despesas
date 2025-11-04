package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

type ExpenseRepository struct {
	db *sqlx.DB
}

func (repo *ExpenseRepository) BulkStore(ctx context.Context, expenses []expense.Expense) error {
	var models []ExpenseModel
	for _, expns := range expenses {
		models = append(models, ToModel(&expns))
	}

	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO expenses (id, name, amount_cents, refund_amount_cents, description, group_id, category_id, split_ratio, split_type, payer_id, receiver_id, created_at, updated_at, deleted_at, version)
    VALUES (:id, :name, :amount_cents, :refund_amount_cents, :description, :group_id, :category_id, :split_ratio, :split_type, :payer_id, :receiver_id, :created_at, :updated_at, :deleted_at, :version)
	`, models); err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

func (repo *ExpenseRepository) GetByGroupDate(ctx context.Context, groupId group.ID, date time.Time) ([]expense.Expense, error) {
	var models []ExpenseModel
	if err := repo.db.SelectContext(ctx, &models, ` 
		SELECT
			id,
			name, 
			amount_cents, 
			refund_amount_cents, 
			description, 
			group_id, 
			category_id, 
			payer_id,   
			receiver_id, 
			split_ratio, 
			split_type, 
			created_at, 
			updated_at, 
			deleted_at, 
			version
		FROM expenses_latest
		WHERE group_id = $1
		AND EXTRACT(MONTH FROM created_at) = $2
		AND EXTRACT(YEAR FROM created_at) = $3
		AND deleted_at IS NULL
		ORDER BY id DESC
  `, groupId.Value, date.Month(), date.Year()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("db.SelectContext: :%w", err)
	}

	var expenses []expense.Expense
	for _, model := range models {
		expenses = append(expenses, *ToEntity(model))
	}

	return expenses, nil
}

func (repo *ExpenseRepository) GetNextID() expense.ID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('expenses_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.QueryRow: %w", err))
	}

	return expense.ID{Value: nextValue}
}

func (repo *ExpenseRepository) GetByID(ctx context.Context, id expense.ID) (*expense.Expense, error) {
	var model ExpenseModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT 
			id, 
			name, 
			amount_cents, 
			refund_amount_cents, 
			description, 
			group_id, 
			category_id, 
			payer_id,   
			receiver_id, 
			split_ratio, 
			split_type, 
			created_at, 
			updated_at, 
			deleted_at, 
			version
		FROM expenses_latest 
		WHERE id = $1 AND deleted_at IS NULL
	`, id.Value).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	return ToEntity(model), nil
}

func (repo *ExpenseRepository) Store(ctx context.Context, entity *expense.Expense) error {
	model := ToModel(entity)

	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO expenses (id, name, amount_cents, refund_amount_cents, description, group_id, category_id, split_ratio, split_type, payer_id, receiver_id, created_at, updated_at, deleted_at, version)
    VALUES (:id, :name, :amount_cents, :refund_amount_cents, :description, :group_id, :category_id, :split_ratio, :split_type, :payer_id, :receiver_id, :created_at, :updated_at, :deleted_at, :version)
	`, &model); err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

func NewExpenseRepository(db *db.Client) expense.Repository {
	return &ExpenseRepository{db: db.Conn()}
}
