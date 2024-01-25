package incomerepo

import (
	"database/sql"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

func toEntity(model IncomeModel) *entity.Income {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &entity.Income{
		Entity: ddd.Entity[entity.IncomeID]{
			ID:        entity.IncomeID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		UserID: entity.UserID{Value: model.UserID},
		Amount: model.Amount,
		Type:   entity.IncomeType(model.Type),
	}
}

func toModel(entity *entity.Income) IncomeModel {
	deletedAt := sql.NullTime{Time: time.Time{}, Valid: false}
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	}

	return IncomeModel{
		ID:        entity.ID.Value,
		UserID:    entity.UserID.Value,
		Amount:    entity.Amount,
		Type:      entity.Type.String(),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
		Version:   entity.Version,
	}
}
