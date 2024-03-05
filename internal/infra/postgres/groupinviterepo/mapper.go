package groupinviterepo

import (
	"database/sql"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

func toEntity(model GroupInviteModel) *entity.GroupInvite {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &entity.GroupInvite{
		Entity: ddd.Entity[entity.GroupInviteID]{
			ID:        entity.GroupInviteID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		GroupID:   entity.GroupID{Value: model.GroupID},
		Token:     model.Token,
		Email:     model.Email,
		ExpiresAt: model.ExpiresAt,
		Status:    entity.GroupInviteStatus(model.Status),
	}
}

func toModel(entity *entity.GroupInvite) GroupInviteModel {
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
