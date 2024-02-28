package categorygrouprepo

import (
	"database/sql"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

func toEntity(model CategoryGroupModel) *entity.CategoryGroup {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &entity.CategoryGroup{
		Entity: ddd.Entity[entity.CategoryGroupID]{
			ID:        entity.CategoryGroupID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name: model.Name,
		Icon: model.Icon,
	}
}

func toModel(entity *entity.CategoryGroup) CategoryGroupModel {
	var deletedAt sql.NullTime
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	} else {
		deletedAt = sql.NullTime{Time: time.Time{}, Valid: false}
	}

	return CategoryGroupModel{
		ID:        entity.ID.Value,
		Name:      entity.Name,
		Icon:      entity.Icon,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
		Version:   entity.Version,
	}
}
