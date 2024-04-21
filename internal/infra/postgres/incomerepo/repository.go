package incomerepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type IncomePGRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db db.Database) repository.IncomeRepository {
	return &IncomePGRepository{db: db.Client()}
}

func (repo *IncomePGRepository) GetNextID() entity.IncomeID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('incomes_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return entity.IncomeID{Value: nextValue}
}

func (repo *IncomePGRepository) GetByID(ctx context.Context, id entity.IncomeID) (*entity.Income, error) {
	var model IncomeModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, user_id, amount_cents, type, created_at, updated_at, deleted_at, version
		FROM incomes WHERE id = $1
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1
	`, id.Value).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	return toEntity(model), nil
}

func (repo *IncomePGRepository) GetUserMonthlyIncomes(ctx context.Context, userID entity.UserID, date *time.Time) ([]entity.Income, error) {
	var incomes []IncomeModel

	d := time.Now()
	if date != nil {
		d = *date
	}

	if err := repo.db.SelectContext(ctx, &incomes, `
		SELECT id, user_id, amount_cents, type, created_at, updated_at, deleted_at, version
		FROM incomes WHERE user_id = $1
		AND EXTRACT(month FROM created_at) = $2
		AND EXTRACT(year FROM created_at) = $3
		AND deleted_at IS NULL
	`, userID.Value, d.Month(), d.Year()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	var entities []entity.Income
	for _, income := range incomes {
		entities = append(entities, *toEntity(income))
	}

	return entities, nil
}

func (repo *IncomePGRepository) Store(ctx context.Context, entity *entity.Income) error {
	model := toModel(entity)
	if err := repo.create(ctx, model); err != nil {
		if err.Error() == "db.Insert: pq: duplicate key value violates unique constraint \"incomes_pkey\"" {
			if err := repo.update(ctx, model); err != nil {
				return fmt.Errorf("repo.update: %w", err)
			}
			return nil
		}
		return fmt.Errorf("repo.create: %w", err)
	}

	return nil
}

func (repo *IncomePGRepository) create(ctx context.Context, model IncomeModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO incomes (id, user_id, amount_cents, type, created_at, updated_at, deleted_at, version)
		VALUES (:id, :user_id, :amount_cents, :type, :created_at, :updated_at, :deleted_at, :version)
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *IncomePGRepository) update(ctx context.Context, model IncomeModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE incomes SET amount_cents = :amount_cents, type = :type, created_at = :created_at, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
		WHERE id = :id and version = :version
	`, model)
	if err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("db.Update: sql: no rows affected")
	}

	return nil
}
