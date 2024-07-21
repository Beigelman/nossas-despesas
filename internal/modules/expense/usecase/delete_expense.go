package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type DeleteExpense func(ctx context.Context, expenseID expense.ID) (*expense.Expense, error)

func NewDeleteExpense(expenseRepo expense.Repository) DeleteExpense {
	return func(ctx context.Context, expenseID expense.ID) (*expense.Expense, error) {
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
