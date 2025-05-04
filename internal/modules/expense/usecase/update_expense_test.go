package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateExpense(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mocks.NewMockuserRepository(t)
	categoryRepo := mocks.NewMockcategoryRepository(t)
	expenseRepo := mocks.NewMockexpenseRepository(t)
	incomeRepo := mocks.NewMockincomeRepository(t)

	grp := group.New(group.Attributes{
		ID:   group.ID{Value: 1},
		Name: "group",
	})

	receiver := user.New(user.Attributes{
		ID:             user.ID{Value: 1},
		Name:           "receiver",
		Email:          "receiver@email.com",
		ProfilePicture: nil,
		GroupID:        &grp.ID,
	})

	payer := user.New(user.Attributes{
		ID:             user.ID{Value: 2},
		Name:           "payer",
		Email:          "payer@email.com",
		ProfilePicture: nil,
		GroupID:        &grp.ID,
	})

	mismatchPayer := user.New(user.Attributes{
		ID:             user.ID{Value: 2},
		Name:           "payer",
		Email:          "payer@email.com",
		ProfilePicture: nil,
		GroupID:        &group.ID{Value: 2},
	})

	catgry := category.New(category.Attributes{
		ID:   category.ID{Value: 1},
		Name: "test category",
		Icon: "1",
	})

	expns, err := expense.New(expense.Attributes{
		ID:          expense.ID{Value: 1},
		Name:        "name",
		Amount:      100,
		Description: "description",
		GroupID:     grp.ID,
		CategoryID:  catgry.ID,
		SplitRatio:  expense.NewEqualSplitRatio(),
		PayerID:     payer.ID,
		ReceiverID:  receiver.ID,
	})
	assert.Nil(t, err)

	updateExpense := usecase.NewUpdateExpense(expenseRepo, userRepo, categoryRepo, incomeRepo)

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(nil, errors.New("test error")).Once()

		p := usecase.UpdateExpenseParams{
			ID:         expns.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &catgry.ID,
			SplitType:  &expense.SplitTypes.Equal,
		}

		expns, err := updateExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "expenseRepo.GetByID: test error")
	})

	t.Run("should return error if expense not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(nil, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:         expns.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &catgry.ID,
			SplitType:  &expense.SplitTypes.Equal,
		}

		expns, err := updateExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "expense not found")
	})

	t.Run("should return error userRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(nil, errors.New("test error")).Once()

		p := usecase.UpdateExpenseParams{
			ID:         expns.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &catgry.ID,
			SplitType:  &expense.SplitTypes.Equal,
		}

		expns, err := updateExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "userRepo.GetByID: test error")
	})

	t.Run("should return error if payer not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(nil, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:         expns.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &catgry.ID,
			SplitType:  &expense.SplitTypes.Equal,
		}

		expns, err := updateExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "payer not found")
	})

	t.Run("should return error if receiver not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(nil, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:         expns.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &catgry.ID,
			SplitType:  &expense.SplitTypes.Equal,
		}

		expns, err := updateExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "receiver not found")
	})

	t.Run("should return error if payer's and receiver's group does not match", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(mismatchPayer, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:         expns.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &catgry.ID,
			SplitType:  &expense.SplitTypes.Equal,
		}

		expns, err := updateExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "group mismatch")
	})

	t.Run("should return error if categoryRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(nil, errors.New("test error")).Once()

		p := usecase.UpdateExpenseParams{
			ID:         expns.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &catgry.ID,
			SplitType:  &expense.SplitTypes.Equal,
		}

		expns, err := updateExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "categoryRepo.GetByID: test error")
	})

	t.Run("should return error if category is not found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(nil, nil).Once()

		p := usecase.UpdateExpenseParams{
			ID:         expns.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &catgry.ID,
			SplitType:  &expense.SplitTypes.Equal,
		}

		expns, err := updateExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "category not found")
	})

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(catgry, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()

		p := usecase.UpdateExpenseParams{
			ID:         expns.ID,
			PayerID:    &payer.ID,
			ReceiverID: &receiver.ID,
			CategoryID: &catgry.ID,
			SplitType:  &expense.SplitTypes.Equal,
		}

		expns, err := updateExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "expenseRepo.Store: test error")
	})

	t.Run("happy path with proporcional split type", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(catgry, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, mock.Anything).Return([]income.Income{{Amount: 60}}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, receiver.ID, mock.Anything).Return([]income.Income{{Amount: 40}}, nil).Once()

		newName := "name 2"
		newAmount := 1000
		newDescription := "description 2"
		p := usecase.UpdateExpenseParams{
			ID:          expns.ID,
			PayerID:     &payer.ID,
			ReceiverID:  &receiver.ID,
			CategoryID:  &catgry.ID,
			SplitType:   &expense.SplitTypes.Proportional,
			Name:        &newName,
			Description: &newDescription,
			Amount:      &newAmount,
		}

		expns, err := updateExpense(ctx, p)
		assert.Equal(t, expense.ID{Value: 1}, expns.ID)
		assert.Equal(t, newName, expns.Name)
		assert.Equal(t, newAmount, expns.Amount)
		assert.Equal(t, newDescription, expns.Description)
		assert.Equal(t, grp.ID, expns.GroupID)
		assert.Equal(t, catgry.ID, expns.CategoryID)
		assert.Equal(t, expense.NewProportionalSplitRatio(60, 40), expns.SplitRatio)
		assert.Equal(t, payer.ID, expns.PayerID)
		assert.Equal(t, receiver.ID, expns.ReceiverID)
		assert.Nil(t, err)
	})

	t.Run("happy path with transfer split type", func(t *testing.T) {
		expenseRepo.EXPECT().GetByID(ctx, expns.ID).Return(expns, nil).Once()
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(catgry, nil).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		newName := "name 2"
		newAmount := 1000
		newDescription := "description 2"
		p := usecase.UpdateExpenseParams{
			ID:          expns.ID,
			PayerID:     &payer.ID,
			ReceiverID:  &receiver.ID,
			CategoryID:  &catgry.ID,
			SplitType:   &expense.SplitTypes.Transfer,
			Name:        &newName,
			Description: &newDescription,
			Amount:      &newAmount,
		}

		expns, err := updateExpense(ctx, p)
		assert.Equal(t, expense.ID{Value: 1}, expns.ID)
		assert.Equal(t, newName, expns.Name)
		assert.Equal(t, newAmount, expns.Amount)
		assert.Equal(t, newDescription, expns.Description)
		assert.Equal(t, grp.ID, expns.GroupID)
		assert.Equal(t, catgry.ID, expns.CategoryID)
		assert.Equal(t, expense.NewTransferRatio(), expns.SplitRatio)
		assert.Equal(t, payer.ID, expns.PayerID)
		assert.Equal(t, receiver.ID, expns.ReceiverID)
		assert.Nil(t, err)
	})
}
