package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"strings"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/jmoiron/sqlx"
)

type GroupInviteRepository struct {
	db *sqlx.DB
}

func NewGroupInviteRepository(db db.Database) group.InviteRepository {
	return &GroupInviteRepository{db: db.Client()}
}

func (repo *GroupInviteRepository) GetNextID() group.InviteID {
	var nextValue int

	if err := repo.db.QueryRowx("SELECT nextval('group_invites_id_seq');").Scan(&nextValue); err != nil {
		panic(fmt.Errorf("db.Select: %w", err))
	}

	return group.InviteID{Value: nextValue}
}

func (repo *GroupInviteRepository) GetByID(ctx context.Context, id group.InviteID) (*group.Invite, error) {
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

	return groupInviteToEntity(model), nil
}

func (repo *GroupInviteRepository) GetByToken(ctx context.Context, token string) (*group.Invite, error) {
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

	return groupInviteToEntity(model), nil
}

func (repo *GroupInviteRepository) GetGroupInvitesByEmail(ctx context.Context, groupID group.ID, email string) ([]group.Invite, error) {
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

	var entities []group.Invite
	for _, model := range models {
		entities = append(entities, *groupInviteToEntity(model))
	}

	return entities, nil
}

func (repo *GroupInviteRepository) Store(ctx context.Context, entity *group.Invite) error {
	model := groupInviteToModel(entity)
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

func (repo *GroupInviteRepository) create(ctx context.Context, model GroupInviteModel) error {
	if _, err := repo.db.NamedExecContext(ctx, `
		INSERT INTO group_invites (id, email, group_id, status, token, expires_at, created_at, updated_at, deleted_at, version)
		VALUES (:id, :email, :group_id, :status, :token, :expires_at, :created_at, :updated_at, :deleted_at, :version)
	`, model); err != nil {
		return fmt.Errorf("db.Insert: %w", err)
	}

	return nil
}

func (repo *GroupInviteRepository) update(ctx context.Context, model GroupInviteModel) error {
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
