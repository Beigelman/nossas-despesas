package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *db.Client) user.Repository {
	return &UserRepository{db: db.Conn()}
}

// GetNextID implements user.UserRepository.
func (repo *UserRepository) GetNextID() user.ID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT NEXTVAL('users_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return user.ID{Value: nextValue}
}

// GetByID implements user.UserRepository.
func (repo *UserRepository) GetByID(ctx context.Context, id user.ID) (*user.User, error) {
	var model UserModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, email, profile_picture, group_id, flags, created_at, updated_at, deleted_at, version
		FROM users WHERE id = $1
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

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var model UserModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, email, profile_picture, group_id, flags, created_at, updated_at, deleted_at, version
		FROM users WHERE email = $1
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1
	`, email).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	return toEntity(model), nil
}

// Store implements user.UserRepository.
func (repo *UserRepository) Store(ctx context.Context, entity *user.User) error {
	model := toModel(entity)
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

func (repo *UserRepository) create(ctx context.Context, model UserModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO users (id, name, email, group_id, profile_picture, flags, created_at, updated_at, deleted_at, version)
    VALUES (:id, :name, :email, :group_id, :profile_picture, :flags, :created_at, :updated_at, :deleted_at, :version)
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *UserRepository) update(ctx context.Context, model UserModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
    UPDATE users SET name = :name, group_id = :group_id, profile_picture = :profile_picture, flags = :flags, updated_at = NOW(), deleted_at = :deleted_at, version = version + 1
		WHERE id = :id AND version = :version
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
