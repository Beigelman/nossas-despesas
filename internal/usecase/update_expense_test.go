package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	mockrepository "github.com/Beigelman/ludaapi/internal/tests/mocks/repository"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateExpense(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mockrepository.NewMockUserRepository(t)
	categoryRepo := mockrepository.NewMockCategoryRepository(t)
	expenseRepo := mockrepository.NewMockExpenseRepository(t)

	group := entity.NewGroup(entity.GroupParams{
		ID:   entity.GroupID{Value: 1},
		Name: "group",
	})

	receiver := entity.NewUser(entity.UserParams{
		ID:             entity.UserID{Value: 1},
		Name:           "receiver",
		Email:          "receiver@email.com",
		ProfilePicture: nil,
		GroupID:        &group.ID,
	})

	payer := entity.NewUser(entity.UserParams{
		ID:             entity.UserID{Value: 2},
		Name:           "payer",
		Email:          "payer@email.com",
		ProfilePicture: nil,
		GroupID:        &group.ID,
	})

	mismatchPayer := entity.NewUser(entity.UserParams{
		ID:             entity.UserID{Value: 2},
		Name:           "payer",
		Email:          "payer@email.com",
		ProfilePicture: nil,
		GroupID:        &entity.GroupID{Value: 2},
	})

	category := entity.NewCategory(entity.CategoryParams{
		ID:   entity.CategoryID{Value: 1},
		Name: "test category",
		Icon: "1",
	})

	expense, err := entity.NewExpense(entity.ExpenseParams{
		ID:          entity.ExpenseID{Value: 1},
		Name:        "name",
		Amount:      100,
		Description: "description",
		GroupID:     group.ID,
		CategoryID:  category.ID,
		SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
		PayerID:     payer.ID,
		ReceiverID:  receiver.ID,
	})
	assert.Nil(t, err)

	updateExpense := usecase.NewUpdateExpense(expenseRepo, userRepo, categoryRepo)

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(nil, errors.New("test error")).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "expenseRepo.GetByID: test error")
	})

	t.Run("should return error if expense not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(nil, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "error=expense not found")
	})

	t.Run("should return error userRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(nil, errors.New("test error")).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "userRepo.GetByID: test error")
	})

	t.Run("should return error if payer not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(nil, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "error=payer not found")
	})

	t.Run("should return error if receiver not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(nil, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "error=receiver not found")
	})

	t.Run("should return error if payer's and receiver's group does not match", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(mismatchPayer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "error=group mismatch")
	})

	t.Run("should return error if categoryRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, category.ID).Return(nil, errors.New("test error")).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "categoryRepo.GetByID: test error")
	})

	t.Run("should return error if category is not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, category.ID).Return(nil, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "error=category not found")
	})

	t.Run("should return error if split ration does not sum 100", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, category.ID).Return(category, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 70, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "error=Unprocessable Entity, internal=entity.Update: expense.Validate: invalid split ratio")
	})

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, category.ID).Return(category, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "expenseRepo.Store: test error")
	})

	t.Run("happy path", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, category.ID).Return(category, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			CategoryID:  category.ID,
			SplitRatio:  entity.SplitRatio{Payer: 70, Receiver: 30},
			Name:        "name 2",
			Amount:      1000,
			Description: "description 2",
		}

		expense, err := updateExpense(ctx, p)
		assert.Equal(t, entity.ExpenseID{Value: 1}, expense.ID)
		assert.Equal(t, "name 2", expense.Name)
		assert.Equal(t, 1000, expense.Amount)
		assert.Equal(t, "description 2", expense.Description)
		assert.Equal(t, group.ID, expense.GroupID)
		assert.Equal(t, category.ID, expense.CategoryID)
		assert.Equal(t, entity.SplitRatio{Payer: 70, Receiver: 30}, expense.SplitRatio)
		assert.Equal(t, payer.ID, expense.PayerID)
		assert.Equal(t, receiver.ID, expense.ReceiverID)
		assert.Nil(t, err)
	})
}
