package postgres

import (
	"database/sql"
	"time"

	"cloud.google.com/go/civil"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
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

func ToScheduledExpenseModel(entity expense.ScheduledExpense) ScheduledExpenseModel {
	var lastGeneratedAt sql.Null[civil.Date]
	if entity.LastGeneratedAt != nil {
		lastGeneratedAt = sql.Null[civil.Date]{V: *entity.LastGeneratedAt, Valid: true}
	}

	return ScheduledExpenseModel{
		ID:              entity.ID.Value,
		Name:            entity.Name,
		AmountCents:     entity.Amount,
		Description:     entity.Description,
		GroupID:         entity.GroupID.Value,
		CategoryID:      entity.CategoryID.Value,
		SplitType:       entity.SplitType.String(),
		PayerID:         entity.PayerID.Value,
		ReceiverID:      entity.ReceiverID.Value,
		FrequencyInDays: entity.FrequencyInDays,
		LastGeneratedAt: lastGeneratedAt,
		IsActive:        entity.IsActive,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
		Version:         entity.Version,
	}
}

func ToScheduledExpenseEntity(model ScheduledExpenseModel) expense.ScheduledExpense {
	var lastGeneratedAt *civil.Date
	if model.LastGeneratedAt.Valid {
		lastGeneratedAt = &model.LastGeneratedAt.V
	}

	return expense.ScheduledExpense{
		Entity: ddd.Entity[expense.ScheduledExpenseID]{
			ID:        expense.ScheduledExpenseID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			Version:   model.Version,
		},
		Name:            model.Name,
		Amount:          model.AmountCents,
		Description:     model.Description,
		GroupID:         group.ID{Value: model.GroupID},
		CategoryID:      category.ID{Value: model.CategoryID},
		SplitType:       expense.SplitType(model.SplitType),
		PayerID:         user.ID{Value: model.PayerID},
		ReceiverID:      user.ID{Value: model.ReceiverID},
		FrequencyInDays: model.FrequencyInDays,
		LastGeneratedAt: lastGeneratedAt,
		IsActive:        model.IsActive,
	}
}
