package expenserepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type ExpensePGRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db db.Database) repository.ExpenseRepository {
	return &ExpensePGRepository{db: db.Client()}
}

// GetNextID implements expense.UserRepository.
func (repo *ExpensePGRepository) GetNextID() entity.ExpenseID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('expenses_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.QueryRow: %w", err))
	}

	return entity.ExpenseID{Value: nextValue}
}

// GetByID implements expense.UserRepository.
func (repo *ExpensePGRepository) GetByID(ctx context.Context, id entity.ExpenseID) (*entity.Expense, error) {
	var model ExpenseModel

	if err := repo.db.QueryRowxContext(ctx, `
		WITH base AS (
			SELECT id, name, amount_cents, refund_amount_cents, description, group_id, category_id, split_ratio, payer_id, receiver_id, created_at, updated_at, deleted_at, version
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

// Store implements expense.UserRepository.
func (repo *ExpensePGRepository) Store(ctx context.Context, entity *entity.Expense) error {
	var model = ToModel(entity)

	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO expenses (id, name, amount_cents, refund_amount_cents,  description, group_id, category_id, split_ratio, payer_id, receiver_id, created_at, updated_at, deleted_at, version)
		VALUES (:id, :name, :amount_cents, :refund_amount_cents, :description, :group_id, :category_id, :split_ratio, :payer_id, :receiver_id, :created_at, :updated_at, :deleted_at, :version)
	`, &model); err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}
