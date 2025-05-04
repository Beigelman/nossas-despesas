package usecase

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
)

type CreateScheduledExpenseInput struct {
	Name            string            `json:"name" validate:"required"`
	Amount          int               `json:"amount" validate:"required"`
	Description     string            `json:"description" validate:"required"`
	GroupID         group.ID          `json:"group_id" validate:"required"`
	CategoryID      category.ID       `json:"category_id" validate:"required"`
	SplitType       expense.SplitType `json:"split_type" validate:"required"`
	PayerID         user.ID           `json:"payer_id" validate:"required"`
	ReceiverID      user.ID           `json:"receiver_id" validate:"required"`
	FrequencyInDays int               `json:"frequency_in_days" validate:"required"`
}

type CreateScheduledExpense func(ctx context.Context, input CreateScheduledExpenseInput) error

func NewCreateScheduledExpense(
	scheduledExpenseRepo expense.ScheduledExpenseRepository,
) CreateScheduledExpense {
	return func(ctx context.Context, input CreateScheduledExpenseInput) error {
		// Criar a despesa agendada
		scheduledExpense, err := expense.NewScheduledExpense(expense.ScheduledExpenseAttributes{
			ID:              scheduledExpenseRepo.GetNextID(),
			Name:            input.Name,
			Amount:          input.Amount,
			Description:     input.Description,
			GroupID:         input.GroupID,
			CategoryID:      input.CategoryID,
			SplitType:       input.SplitType,
			PayerID:         input.PayerID,
			ReceiverID:      input.ReceiverID,
			FrequencyInDays: input.FrequencyInDays,
		})
		if err != nil {
			return fmt.Errorf("failed to create scheduled expense: %w", err)
		}

		// Salvar a despesa agendada
		if err := scheduledExpenseRepo.Store(ctx, scheduledExpense); err != nil {
			return fmt.Errorf("failed to store scheduled expense: %w", err)
		}

		return nil
	}
}
