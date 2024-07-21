package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeleteExpense(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	expenseRepo := mockrepository.NewMockExpenseRepository(t)

	expense, err := expense.New(expense.Attributes{
		ID:          expense.ID{Value: 1},
		Name:        "name",
		Amount:      100,
		Description: "description",
		GroupID:     group.ID{Value: 1},
		CategoryID:  category.CategoryID{Value: 1},
		SplitRatio:  expense.NewEqualSplitRatio(),
		PayerID:     entity.UserID{Value: 1},
		ReceiverID:  entity.UserID{Value: 2},
	})
	assert.Nil(t, err)

	deleteExpense := NewDeleteExpense(expenseRepo)

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
		assert.Equal(t, expense.ExpenseID{Value: 1}, expense.ID)
		assert.NotNil(t, expense.DeletedAt)
		assert.Nil(t, err)
	})
}
