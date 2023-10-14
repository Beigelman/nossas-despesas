package expenserepo

import (
	"database/sql"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

func toEntity(model ExpenseModel) *entity.Expense {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &entity.Expense{
		Entity: ddd.Entity[entity.ExpenseID]{
			ID:        entity.ExpenseID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name:        model.Name,
		Amount:      model.AmountCents,
		Description: model.Description,
		GroupID:     entity.GroupID{Value: model.GroupID},
		CategoryID:  entity.CategoryID{Value: model.CategoryID},
		SplitRatio: entity.SplitRatio{
			Payer:    model.SplitRatio.Payer,
			Receiver: model.SplitRatio.Receiver,
		},
		PayerID:    entity.UserID{Value: model.PayerID},
		ReceiverID: entity.UserID{Value: model.ReceiverID},
	}
}

func toModel(entity *entity.Expense) ExpenseModel {
	var deletedAt sql.NullTime
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	} else {
		deletedAt = sql.NullTime{Time: time.Time{}, Valid: false}
	}

	return ExpenseModel{
		ID:          entity.ID.Value,
		Name:        entity.Name,
		AmountCents: entity.Amount,
		Description: entity.Description,
		GroupID:     entity.GroupID.Value,
		CategoryID:  entity.CategoryID.Value,
		SplitRatio: SplitRatio{
			Payer:    entity.SplitRatio.Payer,
			Receiver: entity.SplitRatio.Receiver,
		},
		PayerID:    entity.PayerID.Value,
		ReceiverID: entity.ReceiverID.Value,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
		DeletedAt:  deletedAt,
		Version:    entity.Version,
	}
}
