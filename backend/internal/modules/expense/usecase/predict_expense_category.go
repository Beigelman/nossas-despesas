package usecase

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/shared/service"
)

type PredictExpenseCategoryInput struct {
	Name   string
	Amount int
}

type PredictExpenseCategory func(ctx context.Context, input PredictExpenseCategoryInput) (int, error)

func NewPredictExpenseCategory(predict service.Predicter) PredictExpenseCategory {
	return func(ctx context.Context, input PredictExpenseCategoryInput) (int, error) {
		categoryID, err := predict.ExpenseCategory(ctx, input.Name, input.Amount)
		if err != nil {
			return 0, fmt.Errorf("predict.ExpenseCategory: %w", err)
		}

		return categoryID, nil
	}
}
