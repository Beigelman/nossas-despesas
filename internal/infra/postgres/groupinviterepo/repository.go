package groupinviterepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type GroupInvitePGRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db db.Database) repository.GroupInviteRepository {
	return &GroupInvitePGRepository{db: db.Client()}
}

func (repo *GroupInvitePGRepository) GetNextID() entity.GroupInviteID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('group_invites_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return entity.GroupInviteID{Value: nextValue}
}

func (repo *GroupInvitePGRepository) GetByID(ctx context.Context, id entity.GroupInviteID) (*entity.GroupInvite, error) {
	var model GroupInviteModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, email, group_id, status, token, expires_at, created_at, updated_at, deleted_at, version
		FROM group_invites WHERE id = $1
		AND deleted_at IS NULL
	`, id.Value).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	return toEntity(model), nil
}

func (repo *GroupInvitePGRepository) GetByToken(ctx context.Context, token string) (*entity.GroupInvite, error) {
	var model GroupInviteModel

	if err := repo.db.QueryRowxContext(ctx, `
		SELECT id, email, group_id, status, token, expires_at, created_at, updated_at, deleted_at, version
		FROM group_invites WHERE token = $1
		AND deleted_at IS NULL
	`, token).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	return toEntity(model), nil
}

func (repo *GroupInvitePGRepository) GetGroupInvitesByEmail(ctx context.Context, groupID entity.GroupID, email string) ([]entity.GroupInvite, error) {
	var models []GroupInviteModel

	if err := repo.db.SelectContext(ctx, &models, `
		SELECT id, email, group_id, status, token, expires_at, created_at, updated_at, deleted_at, version
		FROM group_invites WHERE email = $1
		and group_id = $2
		AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, email, groupID.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	var entities []entity.GroupInvite
	for _, model := range models {
		entities = append(entities, *toEntity(model))
	}

	return entities, nil
}

func (repo *GroupInvitePGRepository) Store(ctx context.Context, entity *entity.GroupInvite) error {
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

func (repo *GroupInvitePGRepository) create(ctx context.Context, model GroupInviteModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO group_invites (id, email, group_id, status, token, expires_at, created_at, updated_at, deleted_at, version)
		VALUES (:id, :email, :group_id, :status, :token, :expires_at, :created_at, :updated_at, :deleted_at, :version)
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *GroupInvitePGRepository) update(ctx context.Context, model GroupInviteModel) error {
	result, err := repo.db.NamedExecContext(ctx, `
		UPDATE group_invites SET status = :status, updated_at = :updated_at, deleted_at = :deleted_at, version = version + 1
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
