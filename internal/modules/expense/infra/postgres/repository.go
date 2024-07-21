package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type ExpenseRepository struct {
	db *sqlx.DB
}

func (repo *ExpenseRepository) BulkStore(ctx context.Context, expenses []expense.Expense) error {
	var models []ExpenseModel
	for _, expense := range expenses {
		models = append(models, ToModel(&expense))
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
      with base as (
				select
    			distinct on (id) id as id,
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
				from expenses
				where group_id = $1
        AND EXTRACT(month FROM created_at) = $2
		    AND EXTRACT(year FROM created_at) = $3
				order by id desc, version desc
			)
		  SELECT * FROM base where deleted_at is null
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
		WITH base AS (
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
			FROM expenses WHERE id = $1
			ORDER BY version DESC
			LIMIT 1
		)
		SELECT * FROM base WHERE deleted_at IS NULL
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

func NewExpenseRepository(db db.Database) expense.Repository {
	return &ExpenseRepository{db: db.Client()}
}
