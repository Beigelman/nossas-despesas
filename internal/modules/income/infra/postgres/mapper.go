package postgres

import (
	"database/sql"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"

	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

func toEntity(model IncomeModel) *income.Income {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &income.Income{
		Entity: ddd.Entity[income.ID]{
			ID:        income.ID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		UserID: user.ID{Value: model.UserID},
		Amount: model.Amount,
		Type:   income.Type(model.Type),
	}
}

func toModel(entity *income.Income) IncomeModel {
	deletedAt := sql.NullTime{Time: time.Time{}, Valid: false}
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	}

	return IncomeModel{
		ID:        entity.ID.Value,
		UserID:    entity.ID.Value,
		Amount:    entity.Amount,
		Type:      entity.Type.String(),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
		Version:   entity.Version,
	}
}
