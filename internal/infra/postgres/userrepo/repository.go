package userrepo

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

type UserPGRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db db.Database) repository.UserRepository {
	return &UserPGRepository{db: db.Client()}
}

// GetNextID implements user.UserRepository.
func (repo *UserPGRepository) GetNextID() entity.UserID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('users_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return entity.UserID{Value: nextValue}
}

// GetByID implements user.UserRepository.
func (repo *UserPGRepository) GetByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	var model UserModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, email, profile_picture, group_id, created_at, updated_at, deleted_at, version
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

func (repo *UserPGRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model UserModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, name, email, profile_picture, group_id, created_at, updated_at, deleted_at, version
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
func (repo *UserPGRepository) Store(ctx context.Context, entity *entity.User) error {
	var model = toModel(entity)
	if err := repo.create(ctx, model); err != nil {
		if err.Error() == "db.Insert: pq: duplicate key value violates unique constraint \"users_pkey\"" {
			if err := repo.update(ctx, model); err != nil {
				return fmt.Errorf("repo.update: %w", err)
			}
			return nil
		}
		return fmt.Errorf("repo.create: %w", err)
	}

	return nil
}

func (repo *UserPGRepository) create(ctx context.Context, model UserModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO users (id, name, email, group_id, profile_picture, created_at, updated_at, deleted_at, version)
		VALUES (:id, :name, :email, :group_id, :profile_picture, :created_at, :updated_at, :deleted_at, :version)
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *UserPGRepository) update(ctx context.Context, model UserModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE users SET name = :name, group_id = :group_id, profile_picture = :profile_picture, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
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
