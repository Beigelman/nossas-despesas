package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type ScheduledExpenseRepository struct {
	db *db.Client
}

func (repo *ScheduledExpenseRepository) GetNextID() expense.ScheduledExpenseID {
	var nextValue int

	conn := repo.db.Conn()

	if err := conn.QueryRowx("SELECT nextval('scheduled_expenses_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.QueryRow: %w", err))
	}

	return expense.ScheduledExpenseID{Value: nextValue}
}

func (repo *ScheduledExpenseRepository) GetByID(ctx context.Context, id expense.ScheduledExpenseID) (*expense.ScheduledExpense, error) {
	var model ScheduledExpenseModel

	conn := repo.db.Conn()

	if err := conn.QueryRowxContext(ctx, `
		SELECT 
			id,
			name,
			amount_cents,
			description,
			group_id,
			category_id,
			split_type,
			payer_id,
			receiver_id,
			frequency_in_days,
			last_generated_at,
			is_active,
			created_at,
			updated_at,
			version
		FROM scheduled_expenses 
		WHERE id = $1
	`, id.Value).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	entity := ToScheduledExpenseEntity(model)

	return &entity, nil
}

func (repo *ScheduledExpenseRepository) GetActiveScheduledExpenses(ctx context.Context) ([]expense.ScheduledExpense, error) {
	conn := repo.db.Conn()
	var models []ScheduledExpenseModel

	if err := conn.SelectContext(ctx, &models, `
		SELECT 
			id,
			name,
			amount_cents,
			description,
			group_id,
			category_id,
			split_type,
			payer_id,
			receiver_id,
			frequency_in_days,
			last_generated_at,
			is_active,
			created_at,
			updated_at,
			version
		FROM scheduled_expenses
		WHERE is_active = true
		AND (
			last_generated_at IS NULL
			OR (last_generated_at + INTERVAL '1 day' * frequency_in_days <= CURRENT_DATE)
		)
	`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	var entities []expense.ScheduledExpense
	for _, model := range models {
		entities = append(entities, ToScheduledExpenseEntity(model))
	}

	return entities, nil
}

func (repo *ScheduledExpenseRepository) Store(ctx context.Context, entity *expense.ScheduledExpense) error {
	return repo.BulkStore(ctx, []expense.ScheduledExpense{*entity})
}

func (repo *ScheduledExpenseRepository) BulkStore(ctx context.Context, entities []expense.ScheduledExpense) error {
	return repo.db.Transaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		for _, entity := range entities {
			model := ToScheduledExpenseModel(entity)

			if _, err := tx.NamedExecContext(ctx, `
				INSERT INTO scheduled_expenses (id, name, amount_cents, description, group_id, category_id, split_type, payer_id, receiver_id, frequency_in_days, last_generated_at, is_active, created_at, updated_at, version) 
				VALUES (:id, :name, :amount_cents, :description, :group_id, :category_id, :split_type, :payer_id, :receiver_id, :frequency_in_days, :last_generated_at, :is_active, :created_at, :updated_at, :version)
				ON CONFLICT (id) DO UPDATE SET
					name = :name,
					amount_cents = :amount_cents,
					description = :description,
					category_id = :category_id,
					split_type = :split_type,
					payer_id = :payer_id,
					receiver_id = :receiver_id,
					frequency_in_days = :frequency_in_days,
					last_generated_at = :last_generated_at,
					is_active = :is_active,
					updated_at = :updated_at,
					version = :version
			`, model); err != nil {
				return fmt.Errorf("db.ExecContext: %w", err)
			}
		}

		return nil
	})
}

func NewScheduledExpenseRepository(db *db.Client) expense.ScheduledExpenseRepository {
	return &ScheduledExpenseRepository{db: db}
}
