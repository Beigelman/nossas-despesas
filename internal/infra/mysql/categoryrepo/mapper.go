package categoryrepo

import (
	"database/sql"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

func toEntity(model CategoryModel) *entity.Category {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &entity.Category{
		Entity: ddd.Entity[entity.CategoryID]{
			ID:        entity.CategoryID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name:            model.Name,
		Icon:            model.Icon,
		GroupCategoryID: entity.CategoryGroupID{Value: model.CategoryGroupID},
	}
}

func toModel(entity *entity.Category) CategoryModel {
	var deletedAt sql.NullTime
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	} else {
		deletedAt = sql.NullTime{Time: time.Time{}, Valid: false}
	}

	return CategoryModel{
		ID:              entity.ID.Value,
		Name:            entity.Name,
		Icon:            entity.Icon,
		CategoryGroupID: entity.GroupCategoryID.Value,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
		DeletedAt:       deletedAt,
		Version:         entity.Version,
	}
}
