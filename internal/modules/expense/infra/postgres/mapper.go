package postgres

import (
	"database/sql"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
	"time"
)

func ToEntity(model ExpenseModel) *expense.Expense {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	var refundAmount *int
	if model.RefundAmountCents.Valid {
		parsedRefundAmount := int(model.RefundAmountCents.Int64)
		refundAmount = &parsedRefundAmount
	}

	return &expense.Expense{
		Entity: ddd.Entity[expense.ID]{
			ID:        expense.ID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name:         model.Name,
		Amount:       model.AmountCents,
		RefundAmount: refundAmount,
		Description:  model.Description,
		GroupID:      group.ID{Value: model.GroupID},
		CategoryID:   category.ID{Value: model.CategoryID},
		SplitRatio: expense.SplitRatio{
			Payer:    model.SplitRatio.Payer,
			Receiver: model.SplitRatio.Receiver,
		},
		SplitType:  expense.SplitType(model.SplitType),
		PayerID:    user.ID{Value: model.PayerID},
		ReceiverID: user.ID{Value: model.ReceiverID},
	}
}

func ToModel(entity *expense.Expense) ExpenseModel {
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
