package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type CategoryGroupRepository struct {
	db *sqlx.DB
}

func NewCategoryGroupRepository(db db.Database) category.GroupRepository {
	return &CategoryGroupRepository{db: db.Client()}
}

func (repo *CategoryGroupRepository) GetNextID() category.GroupID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('category_groups_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.QueryRow: %w", err))
	}

	return category.GroupID{Value: nextValue}
}

func (repo *CategoryGroupRepository) GetByName(ctx context.Context, name string) (*category.Group, error) {
	var model CategoryGroupModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, icon, created_at, updated_at, deleted_at, version
		FROM category_groups WHERE name = $1
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1
	`, name).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.QueryRowContext: %w", err)
	}

	return groupCategoryToEntity(model), nil
}

func (repo *CategoryGroupRepository) GetByID(ctx context.Context, id category.GroupID) (*category.Group, error) {
	var model CategoryGroupModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, icon, created_at, updated_at, deleted_at, version
		FROM category_groups WHERE id = $1
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1
	`, id.Value).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.QueryRowContext: %w", err)
	}

	return groupCategoryToEntity(model), nil
}

func (repo *CategoryGroupRepository) Store(ctx context.Context, entity *category.Group) error {
	model := groupCategoryToModel(entity)
	if err := repo.create(ctx, model); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if err := repo.update(ctx, model); err != nil {
				return fmt.Errorf("repo.update: %w", err)
			}
			return nil
		}
		return fmt.Errorf("repo.create: %w", err)
	}

	return nil
}

func (repo *CategoryGroupRepository) create(ctx context.Context, model CategoryGroupModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO category_groups (id, name, icon, created_at, updated_at, deleted_at, version)
		VALUES (:id, :name, :icon, :created_at, :updated_at, :deleted_at, :version)
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *CategoryGroupRepository) update(ctx context.Context, model CategoryGroupModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE category_groups SET name = :name, icon = :icon, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
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
