package authrepo

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

type AuthPGRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db db.Database) repository.AuthRepository {
	return &AuthPGRepository{db: db.Client()}
}

func (repo *AuthPGRepository) GetNextID() entity.AuthID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('authentications_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return entity.AuthID{Value: nextValue}
}

func (repo *AuthPGRepository) GetByID(ctx context.Context, id entity.AuthID) (*entity.Auth, error) {
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

func (repo *AuthPGRepository) GetByEmail(ctx context.Context, email string, authType entity.AuthType) (*entity.Auth, error) {
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

func (repo *AuthPGRepository) Store(ctx context.Context, entity *entity.Auth) error {
	var model = toModel(entity)
	if err := repo.create(ctx, model); err != nil {
		if err.Error() == "db.Insert: pq: duplicate key value violates unique constraint \"authentications_pkey\"" {
			if err := repo.update(ctx, model); err != nil {
				return fmt.Errorf("repo.update: %w", err)
			}
			return nil
		}
		return fmt.Errorf("repo.create: %w", err)
	}

	return nil
}

func (repo *AuthPGRepository) create(ctx context.Context, model AuthModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO authentications (id, email, password, provider_id, type, created_at, updated_at, deleted_at, version)
		VALUES (:id, :email, :password, :provider_id, :type, :created_at, :updated_at, :deleted_at, :version)
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *AuthPGRepository) update(ctx context.Context, model AuthModel) error {
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
