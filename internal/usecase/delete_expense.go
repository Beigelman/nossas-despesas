package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
)

type DeleteExpense func(ctx context.Context, expenseID entity.ExpenseID) (*entity.Expense, error)

func NewDeleteExpense(expenseRepo repository.ExpenseRepository) DeleteExpense {
	return func(ctx context.Context, expenseID entity.ExpenseID) (*entity.Expense, error) {
		expense, err := expenseRepo.GetByID(ctx, expenseID)
		if err != nil {
			return nil, fmt.Errorf("expenseRepo.GetByID: %w", err)
		}

		if expense == nil {
			return nil, except.NotFoundError("expense not found")
		}

		expense.Delete()

		if err := expenseRepo.Store(ctx, expense); err != nil {
			return nil, fmt.Errorf("expenseRepo.Store: %w", err)
		}

		return expense, nil
	}
}
