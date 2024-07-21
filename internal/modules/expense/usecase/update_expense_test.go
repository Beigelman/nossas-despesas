package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
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
	incomeRepo := mockrepository.NewMockIncomeRepository(t)

	group := group.NewGroup(group.Attributes{
		ID:   group.ID{Value: 1},
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
		GroupID:        &group.GroupID{Value: 2},
	})

	category := category.NewCategory(category.Attributes{
		ID:   category.CategoryID{Value: 1},
		Name: "test category",
		Icon: "1",
	})

	expense, err := expense.New(expense.Attributes{
		ID:          expense.ID{Value: 1},
		Name:        "name",
		Amount:      100,
		Description: "description",
		GroupID:     group.ID,
		CategoryID:  category.ID,
		SplitRatio:  expense.NewEqualSplitRatio(),
		PayerID:     payer.ID,
		ReceiverID:  receiver.ID,
	})
	assert.Nil(t, err)

	updateExpense := NewUpdateExpense(expenseRepo, userRepo, categoryRepo, incomeRepo)

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(nil, errors.New("test error")).Once()

		p := UpdateExpenseParams{
			ID:         expense.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &category.ID,
			SplitType:  &expense.SpliteTypes.Equal,
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "expenseRepo.GetByID: test error")
	})

	t.Run("should return error if expense not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(nil, nil).Once()

		p := UpdateExpenseParams{
			ID:         expense.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &category.ID,
			SplitType:  &expense.SpliteTypes.Equal,
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "expense not found")
	})

	t.Run("should return error userRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(nil, errors.New("test error")).Once()

		p := UpdateExpenseParams{
			ID:         expense.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &category.ID,
			SplitType:  &expense.SpliteTypes.Equal,
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "userRepo.GetByID: test error")
	})

	t.Run("should return error if payer not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(nil, nil).Once()

		p := UpdateExpenseParams{
			ID:         expense.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &category.ID,
			SplitType:  &expense.SpliteTypes.Equal,
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "payer not found")
	})

	t.Run("should return error if receiver not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(nil, nil).Once()

		p := UpdateExpenseParams{
			ID:         expense.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &category.ID,
			SplitType:  &expense.SpliteTypes.Equal,
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "receiver not found")
	})

	t.Run("should return error if payer's and receiver's group does not match", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(mismatchPayer, nil).Once()

		p := UpdateExpenseParams{
			ID:         expense.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &category.ID,
			SplitType:  &expense.SpliteTypes.Equal,
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "group mismatch")
	})

	t.Run("should return error if categoryRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, category.ID).Return(nil, errors.New("test error")).Once()

		p := UpdateExpenseParams{
			ID:         expense.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &category.ID,
			SplitType:  &expense.SpliteTypes.Equal,
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

		p := UpdateExpenseParams{
			ID:         expense.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &category.ID,
			SplitType:  &expense.SpliteTypes.Equal,
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "category not found")
	})

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, category.ID).Return(category, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()

		p := UpdateExpenseParams{
			ID:         expense.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &category.ID,
			SplitType:  &expense.SpliteTypes.Equal,
		}

		expense, err := updateExpense(ctx, p)
		assert.Nil(t, expense)
		assert.EqualError(t, err, "expenseRepo.Store: test error")
	})

	t.Run("happy path with proporcional split type", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, category.ID).Return(category, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, mock.Anything).Return([]income.Income{{Amount: 60}}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, receiver.ID, mock.Anything).Return([]income.Income{{Amount: 40}}, nil).Once()

		newName := "name 2"
		newAmount := 1000
		newDescription := "description 2"
		p := UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     &payer.ID,
			ReceiverID:  &receiver.ID,
			CategoryID:  &category.ID,
			SplitType:   &expense.SpliteTypes.Proportional,
			Name:        &newName,
			Description: &newDescription,
			Amount:      &newAmount,
		}

		expense, err := updateExpense(ctx, p)
		assert.Equal(t, expense.ExpenseID{Value: 1}, expense.ID)
		assert.Equal(t, newName, expense.Name)
		assert.Equal(t, newAmount, expense.Amount)
		assert.Equal(t, newDescription, expense.Description)
		assert.Equal(t, group.ID, expense.GroupID)
		assert.Equal(t, category.ID, expense.CategoryID)
		assert.Equal(t, expense.NewProportionalSplitRatio(60, 40), expense.SplitRatio)
		assert.Equal(t, payer.ID, expense.PayerID)
		assert.Equal(t, receiver.ID, expense.ReceiverID)
		assert.Nil(t, err)
	})

	t.Run("happy path with transfer split type", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expense.ID).Return(expense, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, category.ID).Return(category, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		newName := "name 2"
		newAmount := 1000
		newDescription := "description 2"
		p := UpdateExpenseParams{
			ID:          expense.ID,
			PayerID:     &payer.ID,
			ReceiverID:  &receiver.ID,
			CategoryID:  &category.ID,
			SplitType:   &expense.SpliteTypes.Transfer,
			Name:        &newName,
			Description: &newDescription,
			Amount:      &newAmount,
		}

		expense, err := updateExpense(ctx, p)
		assert.Equal(t, expense.ExpenseID{Value: 1}, expense.ID)
		assert.Equal(t, newName, expense.Name)
		assert.Equal(t, newAmount, expense.Amount)
		assert.Equal(t, newDescription, expense.Description)
		assert.Equal(t, group.ID, expense.GroupID)
		assert.Equal(t, category.ID, expense.CategoryID)
		assert.Equal(t, expense.NewTransferRatio(), expense.SplitRatio)
		assert.Equal(t, payer.ID, expense.PayerID)
		assert.Equal(t, receiver.ID, expense.ReceiverID)
		assert.Nil(t, err)
	})
}
