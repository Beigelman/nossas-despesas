package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeleteExpense(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	expenseRepo := mockrepository.NewMockExpenseRepository(t)

	expns, err := expense.New(expense.Attributes{
		ID:          expense.ID{Value: 1},
		Name:        "name",
		Amount:      100,
		Description: "description",
		GroupID:     group.ID{Value: 1},
		CategoryID:  category.ID{Value: 1},
		SplitRatio:  expense.NewEqualSplitRatio(),
		PayerID:     user.ID{Value: 1},
		ReceiverID:  user.ID{Value: 2},
	})
	assert.Nil(t, err)

	deleteExpense := usecase.NewDeleteExpense(expenseRepo)

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(nil, errors.New("test error")).Once()

		delExpense, err := deleteExpense(ctx, expns.ID)
		assert.Nil(t, delExpense)
		assert.EqualError(t, err, "expenseRepo.GetByID: test error")
	})

	t.Run("should return error if expense not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(nil, nil).Once()

		delExpense, err := deleteExpense(ctx, expns.ID)
		assert.Nil(t, delExpense)
		assert.EqualError(t, err, "expense not found")
	})

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()

		delExpense, err := deleteExpense(ctx, expns.ID)
		assert.Nil(t, delExpense)
		assert.EqualError(t, err, "expenseRepo.Store: test error")
	})

	t.Run("happy path", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		delExpense, err := deleteExpense(ctx, expns.ID)
		assert.Equal(t, expense.ID{Value: 1}, delExpense.ID)
		assert.NotNil(t, delExpense.DeletedAt)
		assert.Nil(t, err)
	})
}
