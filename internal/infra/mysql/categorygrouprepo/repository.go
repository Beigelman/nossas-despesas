package categorygrouprepo

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

type MySqlCategoryGroupRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db db.Database) repository.CategoryGroupRepository {
	return &MySqlCategoryGroupRepository{db: db.Client()}
}

func (repo *MySqlCategoryGroupRepository) GetNextID() entity.CategoryGroupID {
	var nextValue int
	if _, err := repo.db.Exec("SET information_schema_stats_expiry = 0;"); err != nil {
		panic(fmt.Errorf("db.Exec: %w", err))
	}

	if err := repo.db.QueryRowx(`
		SELECT AUTO_INCREMENT
		FROM information_schema.tables
		WHERE table_name = 'category_groups'
		AND table_schema = DATABASE();
	`).Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return entity.CategoryGroupID{Value: nextValue}
}

func (repo *MySqlCategoryGroupRepository) GetByName(ctx context.Context, name string) (*entity.CategoryGroup, error) {
	var model CategoryGroupModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, icon, created_at, updated_at, deleted_at, version
		FROM category_groups WHERE name = ?
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

func (repo *MySqlCategoryGroupRepository) GetByID(ctx context.Context, id entity.CategoryGroupID) (*entity.CategoryGroup, error) {
	var model CategoryGroupModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, icon, created_at, updated_at, deleted_at, version
		FROM category_groups WHERE id = ?
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

func (repo *MySqlCategoryGroupRepository) Store(ctx context.Context, entity *entity.CategoryGroup) error {
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

func (repo *MySqlCategoryGroupRepository) create(ctx context.Context, model CategoryGroupModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO category_groups (id, name, icon, created_at, updated_at, deleted_at, version)
		VALUES (:id, :name, :icon, :created_at, :updated_at, :deleted_at, :version);
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *MySqlCategoryGroupRepository) update(ctx context.Context, model CategoryGroupModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE category_groups SET name = :name, icon = :icon, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
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
