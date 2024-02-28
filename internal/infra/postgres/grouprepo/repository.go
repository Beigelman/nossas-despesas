package grouprepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type GroupPGRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db db.Database) repository.GroupRepository {
	return &GroupPGRepository{db: db.Client()}
}

// GetNextID implements group.UserRepository.
func (repo *GroupPGRepository) GetNextID() entity.GroupID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('groups_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return entity.GroupID{Value: nextValue}
}

// GetByID implements group.UserRepository.
func (repo *GroupPGRepository) GetByID(ctx context.Context, id entity.GroupID) (*entity.Group, error) {
	var model GroupModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, created_at, updated_at, deleted_at, version
		FROM groups WHERE id = $1
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

// GetByName implements group.UserRepository.
func (repo *GroupPGRepository) GetByName(ctx context.Context, name string) (*entity.Group, error) {
	var model GroupModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, created_at, updated_at, deleted_at, version
		FROM groups WHERE name = $1
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1
	`, name).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	return toEntity(model), nil
}

// Store implements group.UserRepository.
func (repo *GroupPGRepository) Store(ctx context.Context, entity *entity.Group) error {
	var model = toModel(entity)

	if err := repo.create(ctx, model); err != nil {
		if err.Error() == "db.Insert: pq: duplicate key value violates unique constraint \"groups_pkey\"" {
			if err := repo.update(ctx, model); err != nil {
				return fmt.Errorf("repo.update: %w", err)
			}
			return nil
		}
		return fmt.Errorf("repo.create: %w", err)
	}

	return nil
}

func (repo *GroupPGRepository) create(ctx context.Context, model GroupModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO groups (id, name, created_at, updated_at, deleted_at, version)
		VALUES (:id, :name, :created_at, :updated_at, :deleted_at, :version)
	`, &model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *GroupPGRepository) update(ctx context.Context, model GroupModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE groups SET name = :name, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
		WHERE id = :id AND version = :version
	`, &model)
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
