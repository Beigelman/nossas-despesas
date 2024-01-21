package authrepo

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

type MySqlAuthRepository struct {
	db *sqlx.DB
}

func NewMySqlRepository(db db.Database) repository.AuthRepository {
	return &MySqlAuthRepository{db: db.Client()}
}

func (repo *MySqlAuthRepository) GetNextID() entity.AuthID {
	var nextValue int
	if _, err := repo.db.Exec("SET information_schema_stats_expiry = 0;"); err != nil {
		panic(fmt.Errorf("db.Exec: %w", err))
	}

	if err := repo.db.QueryRowx(`
		SELECT AUTO_INCREMENT
		FROM information_schema.tables
		WHERE table_name = 'authentications'
		AND table_schema = DATABASE();
	`).Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return entity.AuthID{Value: nextValue}
}

func (repo *MySqlAuthRepository) GetByID(ctx context.Context, id entity.AuthID) (*entity.Auth, error) {
	var model AuthModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, email, password, provider_id, type, created_at, updated_at, deleted_at, version
		FROM authentications WHERE id = ?
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1;
	`, id.Value).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	return toEntity(model), nil
}

func (repo *MySqlAuthRepository) GetByEmail(ctx context.Context, email string, authType entity.AuthType) (*entity.Auth, error) {
	var model AuthModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, email, password, provider_id, type, created_at, updated_at, deleted_at, version
		FROM authentications WHERE email = ? and type = ?
		AND deleted_at IS NULL
		ORDER BY version DESC
		LIMIT 1;
	`, email, string(authType)).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	return toEntity(model), nil
}

func (repo *MySqlAuthRepository) Store(ctx context.Context, entity *entity.Auth) error {
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

func (repo *MySqlAuthRepository) create(ctx context.Context, model AuthModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO authentications (id, email, password, provider_id, type, created_at, updated_at, deleted_at, version)
		VALUES (:id, :email, :password, :provider_id, :type, :created_at, :updated_at, :deleted_at, :version);
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *MySqlAuthRepository) update(ctx context.Context, model AuthModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE authentications SET password = :password, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
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
