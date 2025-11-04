package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
)

func TestPredictExpenseCategory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	predictService := mocks.NewMockservicePredicter(t)

	predictExpenseCategory := usecase.NewPredictExpenseCategory(predictService)

	input := usecase.PredictExpenseCategoryInput{
		Name:   "Farmácia",
		Amount: 5000,
	}

	t.Run("should return error if predict service fails", func(t *testing.T) {
		predictService.EXPECT().ExpenseCategory(ctx, input.Name, input.Amount).Return(0, errors.New("test error")).Once()

		categoryID, err := predictExpenseCategory(ctx, input)
		assert.Equal(t, 0, categoryID)
		assert.EqualError(t, err, "predict.ExpenseCategory: test error")
	})

	t.Run("should return category ID on success", func(t *testing.T) {
		expectedCategoryID := 1
		predictService.EXPECT().ExpenseCategory(ctx, input.Name, input.Amount).Return(expectedCategoryID, nil).Once()

		categoryID, err := predictExpenseCategory(ctx, input)
		assert.Equal(t, expectedCategoryID, categoryID)
		assert.Nil(t, err)
	})

	t.Run("should handle different expense names and amounts", func(t *testing.T) {
		testCases := []struct {
			name   string
			amount int
			result int
		}{
			{"Gasolina", 25000, 2},
			{"Aluguel", 350000, 3},
			{"Supermercado", 15000, 4},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				testInput := usecase.PredictExpenseCategoryInput{
					Name:   tc.name,
					Amount: tc.amount,
				}
				predictService.EXPECT().ExpenseCategory(ctx, tc.name, tc.amount).Return(tc.result, nil).Once()

				categoryID, err := predictExpenseCategory(ctx, testInput)
				assert.Equal(t, tc.result, categoryID)
				assert.Nil(t, err)
			})
		}
	})

	t.Run("should pass context correctly", func(t *testing.T) {
		predictService.EXPECT().ExpenseCategory(mock.Anything, input.Name, input.Amount).
			Run(func(ctx context.Context, name string, amount int) {
				assert.NotNil(t, ctx)
			}).
			Return(1, nil).Once()

		categoryID, err := predictExpenseCategory(ctx, input)
		assert.Equal(t, 1, categoryID)
		assert.Nil(t, err)
	})
}
