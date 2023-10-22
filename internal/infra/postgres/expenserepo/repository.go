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

type PGRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db db.Database) repository.ExpenseRepository {
	return &PGRepository{db: db.Client()}
}

// GetNextID implements expense.UserRepository.
func (repo *PGRepository) GetNextID() entity.ExpenseID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('expenses_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.QueryRow: %w", err))
	}

	return entity.ExpenseID{Value: nextValue}
}

// GetByID implements expense.UserRepository.
func (repo *PGRepository) GetByID(ctx context.Context, id entity.ExpenseID) (*entity.Expense, error) {
	var model ExpenseModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, amount_cents, description, group_id, category_id, split_ratio, payer_id, receiver_id, created_at, updated_at, deleted_at, version
		FROM expenses WHERE id = $1
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1
	`, id.Value).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	return toEntity(model), nil
}

// Store implements expense.UserRepository.
func (repo *PGRepository) Store(ctx context.Context, entity *entity.Expense) error {
	var model = toModel(entity)

	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO expenses (id, name, amount_cents, description, group_id, category_id, split_ratio, payer_id, receiver_id, created_at, updated_at, deleted_at, version)
		VALUES (:id, :name, :amount_cents, :description, :group_id, :category_id, :split_ratio, :payer_id, :receiver_id, :created_at, :updated_at, :deleted_at, :version)
	`, &model); err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}
