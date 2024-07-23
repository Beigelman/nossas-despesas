package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"strings"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db db.Database) category.Repository {
	return &CategoryRepository{db: db.Client()}
}

func (repo *CategoryRepository) GetNextID() category.ID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('categories_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.QueryRow: %w", err))
	}

	return category.ID{Value: nextValue}
}

func (repo *CategoryRepository) GetByName(ctx context.Context, name string) (*category.Category, error) {
	var model CategoryModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, icon, category_group_id,created_at, updated_at, deleted_at, version
		FROM categories WHERE name = $1
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1
	`, name).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.QueryRowContext: %w", err)
	}

	return categoryToEntity(model), nil
}

func (repo *CategoryRepository) GetByID(ctx context.Context, id category.ID) (*category.Category, error) {
	var model CategoryModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, icon, category_group_id, created_at, updated_at, deleted_at, version
		FROM categories WHERE id = $1
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1
	`, id.Value).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.QueryRowContext: %w", err)
	}

	return categoryToEntity(model), nil
}

func (repo *CategoryRepository) Store(ctx context.Context, entity *category.Category) error {
	model := categoryToModel(entity)
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

func (repo *CategoryRepository) create(ctx context.Context, model CategoryModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO categories (id, name, icon, category_group_id, created_at, updated_at, deleted_at, version)
		VALUES (:id, :name, :icon, :category_group_id,  :created_at, :updated_at, :deleted_at, :version)
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *CategoryRepository) update(ctx context.Context, model CategoryModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE categories SET name = :name, icon = :icon, category_group_id = :category_group_id, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
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
