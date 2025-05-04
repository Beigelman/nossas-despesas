package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db db.Database) auth.Repository {
	return &AuthRepository{db: db.Client()}
}

func (repo *AuthRepository) GetNextID() auth.ID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('authentications_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return auth.ID{Value: nextValue}
}

func (repo *AuthRepository) GetByID(ctx context.Context, id auth.ID) (*auth.Auth, error) {
	var model AuthModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, email, password, provider_id, type, created_at, updated_at, deleted_at, version
		FROM authentications WHERE id = $1
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

func (repo *AuthRepository) GetByEmail(ctx context.Context, email string, authType auth.Type) (*auth.Auth, error) {
	var model AuthModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, email, password, provider_id, type, created_at, updated_at, deleted_at, version
		FROM authentications WHERE email = $1 and type = $2
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1
	`, email, string(authType)).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	return toEntity(model), nil
}

func (repo *AuthRepository) Store(ctx context.Context, entity *auth.Auth) error {
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

func (repo *AuthRepository) create(ctx context.Context, model AuthModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO authentications (id, email, password, provider_id, type, created_at, updated_at, deleted_at, version)
		VALUES (:id, :email, :password, :provider_id, :type, :created_at, :updated_at, :deleted_at, :version)
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *AuthRepository) update(ctx context.Context, model AuthModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE authentications SET password = :password, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
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
