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

type MySqlExpensesRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db db.Database) repository.ExpenseRepository {
	return &MySqlExpensesRepository{db: db.Client()}
}

// GetNextID implements expense.UserRepository.
func (repo *MySqlExpensesRepository) GetNextID() entity.ExpenseID {
	var nextValue int
	if _, err := repo.db.Exec("SET information_schema_stats_expiry = 0;"); err != nil {
		panic(fmt.Errorf("db.Exec: %w", err))
	}

	if err := repo.db.QueryRowx(`
		SELECT AUTO_INCREMENT
		FROM information_schema.tables
		WHERE table_name = 'expenses'
		AND table_schema = DATABASE();
	`).Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return entity.ExpenseID{Value: nextValue}
}

// GetByID implements expense.UserRepository.
func (repo *MySqlExpensesRepository) GetByID(ctx context.Context, id entity.ExpenseID) (*entity.Expense, error) {
	var model ExpenseModel

	if err := repo.db.QueryRowxContext(ctx, `
		WITH base AS (
			SELECT id, name, amount_cents, description, group_id, category_id, split_ratio, payer_id, receiver_id, created_at, updated_at, deleted_at, version
			FROM expenses WHERE id = ?
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

	return toEntity(model), nil
}

// Store implements expense.UserRepository.
func (repo *MySqlExpensesRepository) Store(ctx context.Context, entity *entity.Expense) error {
	var model = toModel(entity)

	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, split_ratio, payer_id, receiver_id, created_at, updated_at, deleted_at, version)
		VALUES (:id, :name, :amount_cents, :description, :group_id, :category_id, :split_ratio, :payer_id, :receiver_id, :created_at, :updated_at, :deleted_at, :version)
	`, &model); err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}
