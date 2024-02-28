package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	vo "github.com/Beigelman/nossas-despesas/internal/domain/valueobject"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
	"github.com/Beigelman/nossas-despesas/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeleteExpense(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	expenseRepo := mockrepository.NewMockExpenseRepository(t)

	expense, err := entity.NewExpense(entity.ExpenseParams{
		ID:          entity.ExpenseID{Value: 1},
		Name:        "name",
		Amount:      100,
		Description: "description",
		GroupID:     entity.GroupID{Value: 1},
		CategoryID:  entity.CategoryID{Value: 1},
		SplitRatio:  vo.NewEqualSplitRatio(),
		PayerID:     entity.UserID{Value: 1},
		ReceiverID:  entity.UserID{Value: 2},
	})
	assert.Nil(t, err)

	deleteExpense := usecase.NewDeleteExpense(expenseRepo)

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(nil, errors.New("test error")).Once()

		expense, err := deleteExpense(ctx, expense.ID)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "expenseRepo.GetByID: test error")
	})

	t.Run("should return error if expense not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(nil, nil).Once()

		expense, err := deleteExpense(ctx, expense.ID)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "expense not found")
	})

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()

		expense, err := deleteExpense(ctx, expense.ID)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "expenseRepo.Store: test error")
	})

	t.Run("happy path", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		expense, err := deleteExpense(ctx, expense.ID)
		assert.Equal(t, entity.ExpenseID{Value: 1}, expense.ID)
		assert.NotNil(t, expense.DeletedAt)
		assert.Nil(t, err)
	})
}
