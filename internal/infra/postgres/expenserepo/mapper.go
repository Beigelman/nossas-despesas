package expenserepo

import (
	"database/sql"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	vo "github.com/Beigelman/nossas-despesas/internal/domain/valueobject"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

func ToEntity(model ExpenseModel) *entity.Expense {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	var refundAmount *int
	if model.RefundAmountCents.Valid {
		parsedRefundAmount := int(model.RefundAmountCents.Int64)
		refundAmount = &parsedRefundAmount
	}

	return &entity.Expense{
		Entity: ddd.Entity[entity.ExpenseID]{
			ID:        entity.ExpenseID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name:         model.Name,
		Amount:       model.AmountCents,
		RefundAmount: refundAmount,
		Description:  model.Description,
		GroupID:      entity.GroupID{Value: model.GroupID},
		CategoryID:   entity.CategoryID{Value: model.CategoryID},
		SplitRatio: vo.SplitRatio{
			Payer:    model.SplitRatio.Payer,
			Receiver: model.SplitRatio.Receiver,
		},
		SplitType:  vo.SplitType(model.SplitType),
		PayerID:    entity.UserID{Value: model.PayerID},
		ReceiverID: entity.UserID{Value: model.ReceiverID},
	}
}

func ToModel(entity *entity.Expense) ExpenseModel {
	var deletedAt sql.NullTime
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	}

	var refundAmount sql.NullInt64
	if entity.RefundAmount != nil {
		refundAmount = sql.NullInt64{Int64: int64(*entity.RefundAmount), Valid: true}
	}

	return ExpenseModel{
		ID:                entity.ID.Value,
		Name:              entity.Name,
		AmountCents:       entity.Amount,
		RefundAmountCents: refundAmount,
		Description:       entity.Description,
		GroupID:           entity.GroupID.Value,
		CategoryID:        entity.CategoryID.Value,
		SplitRatio: SplitRatio{
			Payer:    entity.SplitRatio.Payer,
			Receiver: entity.SplitRatio.Receiver,
		},
		SplitType:  entity.SplitType.String(),
		PayerID:    entity.PayerID.Value,
		ReceiverID: entity.ReceiverID.Value,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
		DeletedAt:  deletedAt,
		Version:    entity.Version,
	}
}
