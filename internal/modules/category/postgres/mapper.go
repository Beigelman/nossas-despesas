package postgres

import (
	"database/sql"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

func categoryToEntity(model CategoryModel) *category.Category {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &category.Category{
		Entity: ddd.Entity[category.ID]{
			ID:        category.ID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name:            model.Name,
		Icon:            model.Icon,
		GroupCategoryID: category.GroupID{Value: model.CategoryGroupID},
	}
}

func categoryToModel(entity *category.Category) CategoryModel {
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

func groupCategoryToEntity(model CategoryGroupModel) *category.Group {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &category.Group{
		Entity: ddd.Entity[category.GroupID]{
			ID:        category.GroupID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name: model.Name,
		Icon: model.Icon,
	}
}

func groupCategoryToModel(entity *category.Group) CategoryGroupModel {
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
