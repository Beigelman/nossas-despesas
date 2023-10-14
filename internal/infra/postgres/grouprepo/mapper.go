package grouprepo

import (
	"database/sql"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

func toEntity(model GroupModel) *entity.Group {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &entity.Group{
		Entity: ddd.Entity[entity.GroupID]{
			ID:        entity.GroupID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name: model.Name,
	}
}

func toModel(entity *entity.Group) GroupModel {
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
