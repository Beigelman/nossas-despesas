package categoryrepo

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

type MySqlCategoryRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db db.Database) repository.CategoryRepository {
	return &MySqlCategoryRepository{db: db.Client()}
}

func (repo *MySqlCategoryRepository) GetNextID() entity.CategoryID {
	var nextValue int
	if _, err := repo.db.Exec("SET information_schema_stats_expiry = 0;"); err != nil {
		panic(fmt.Errorf("db.Exec: %w", err))
	}

	if err := repo.db.QueryRowx(`
		SELECT AUTO_INCREMENT
		FROM information_schema.tables
		WHERE table_name = 'categories'
		AND table_schema = DATABASE();
	`).Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return entity.CategoryID{Value: nextValue}
}

func (repo *MySqlCategoryRepository) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	var model CategoryModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, icon, category_group_id,created_at, updated_at, deleted_at, version
		FROM categories WHERE name = ?
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1;
	`, name).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.QueryRowContext: %w", err)
	}

	return toEntity(model), nil
}

func (repo *MySqlCategoryRepository) GetByID(ctx context.Context, id entity.CategoryID) (*entity.Category, error) {
	var model CategoryModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, icon, category_group_id, created_at, updated_at, deleted_at, version
		FROM categories WHERE id = ?
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1;
	`, id.Value).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.QueryRowContext: %w", err)
	}

	return toEntity(model), nil
}

func (repo *MySqlCategoryRepository) Store(ctx context.Context, entity *entity.Category) error {
	var model = toModel(entity)
	if err := repo.create(ctx, model); err != nil {
		if err.Error() == "db.Insert: pq: duplicate key value violates unique constraint \"categories_pkey\"" {
			if err := repo.update(ctx, model); err != nil {
				return fmt.Errorf("repo.update: %w", err)
			}
			return nil
		}
		return fmt.Errorf("repo.create: %w", err)
	}

	return nil
}

func (repo *MySqlCategoryRepository) create(ctx context.Context, model CategoryModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO categories (id, name, icon, category_group_id, created_at, updated_at, deleted_at, version)
		VALUES (:id, :name, :icon, :category_group_id,  :created_at, :updated_at, :deleted_at, :version);
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *MySqlCategoryRepository) update(ctx context.Context, model CategoryModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE categories SET name = :name, icon = :icon, category_group_id = :category_group_id, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
		WHERE id = :id and version = :version;
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
