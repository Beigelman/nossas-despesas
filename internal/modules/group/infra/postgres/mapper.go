package postgres

import (
	"database/sql"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

func groupToEntity(model GroupModel) *group.Group {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &group.Group{
		Entity: ddd.Entity[group.ID]{
			ID:        group.ID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name: model.Name,
	}
}

func groupToModel(entity *group.Group) GroupModel {
	var deletedAt sql.NullTime
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	} else {
		deletedAt = sql.NullTime{Time: time.Time{}, Valid: false}
	}

	return GroupModel{
		ID:        entity.ID.Value,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
		Version:   entity.Version,
	}
}

func groupInviteToEntity(model GroupInviteModel) *group.Invite {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &group.Invite{
		Entity: ddd.Entity[group.InviteID]{
			ID:        group.InviteID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		GroupID:   group.ID{Value: model.GroupID},
		Token:     model.Token,
		Email:     model.Email,
		ExpiresAt: model.ExpiresAt,
		Status:    group.InviteStatus(model.Status),
	}
}

func groupInviteToModel(entity *group.Invite) GroupInviteModel {
	deletedAt := sql.NullTime{Time: time.Time{}, Valid: false}
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	}

	return GroupInviteModel{
		ID:        entity.ID.Value,
		GroupID:   entity.GroupID.Value,
		Email:     entity.Email,
		Token:     entity.Token,
		Status:    string(entity.Status),
		ExpiresAt: entity.ExpiresAt,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
		Version:   entity.Version,
	}
}
