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

func TestCreateExpense(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mocks.NewMockuserRepository(t)
	groupRepo := mocks.NewMockgroupRepository(t)
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

	catgry := category.New(category.Attributes{
		ID:   category.ID{Value: 1},
		Name: "test category",
		Icon: "1",
	})

	createExpense := usecase.NewCreateExpense(expenseRepo, userRepo, groupRepo, categoryRepo, incomeRepo)

	t.Run("should return error userRepo fails", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(nil, errors.New("test error")).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "userRepo.GetByID: test error")
	})

	t.Run("should return error if payer not found", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(nil, nil).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "payer not found")
	})

	t.Run("should return error if receiver not found", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(nil, nil).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "receiver not found")
	})

	t.Run("should return error if groupRepo fails", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(nil, errors.New("test error")).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "groupRepo.GetByID: test error")
	})

	t.Run("should return error if group not found", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(nil, nil).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "group not found")
	})

	t.Run("should return error if payer's and receiver's group does not match", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, group.ID{Value: 3}).Return(group.New(group.Attributes{ID: group.ID{Value: 3}}), nil).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     group.ID{Value: 3},
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "group mismatch")
	})

	t.Run("should return error if categoryRepo fails", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(nil, errors.New("test error")).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "categoryRepo.GetByID: test error")
	})

	t.Run("should return error if category is not found", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(nil, nil).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "category not found")
	})

	t.Run("should return error if expenseRepo fails", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(catgry, nil).Once()
		expenseRepo.EXPECT().GetNextID().Return(expense.ID{Value: 1}).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Nil(t, expns)
		assert.EqualError(t, err, "expenseRepo.Store: test error")
	})

	t.Run("happy path with equal split ratio", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(catgry, nil).Once()
		expenseRepo.EXPECT().GetNextID().Return(expense.ID{Value: 1}).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "equal",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Equal(t, expense.ID{Value: 1}, expns.ID)
		assert.Equal(t, "name", expns.Name)
		assert.Equal(t, 100, expns.Amount)
		assert.Equal(t, "description", expns.Description)
		assert.Equal(t, grp.ID, expns.GroupID)
		assert.Equal(t, catgry.ID, expns.CategoryID)
		assert.Equal(t, expense.NewEqualSplitRatio(), expns.SplitRatio)
		assert.Equal(t, payer.ID, expns.PayerID)
		assert.Equal(t, receiver.ID, expns.ReceiverID)
		assert.Nil(t, err)
	})

	t.Run("happy path with proportional split ratio", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(catgry, nil).Once()
		expenseRepo.EXPECT().GetNextID().Return(expense.ID{Value: 1}).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, mock.Anything).Return([]income.Income{{Amount: 40}}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, receiver.ID, mock.Anything).Return([]income.Income{{Amount: 60}}, nil).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "proportional",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Equal(t, expense.ID{Value: 1}, expns.ID)
		assert.Equal(t, "name", expns.Name)
		assert.Equal(t, 100, expns.Amount)
		assert.Equal(t, "description", expns.Description)
		assert.Equal(t, grp.ID, expns.GroupID)
		assert.Equal(t, catgry.ID, expns.CategoryID)
		assert.Equal(t, expense.SplitRatio{
			Payer:    40,
			Receiver: 60,
		}, expns.SplitRatio)
		assert.Equal(t, payer.ID, expns.PayerID)
		assert.Equal(t, receiver.ID, expns.ReceiverID)
		assert.Nil(t, err)
	})

	t.Run("happy path with transfer split ratio", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, payer.ID).Return(payer, nil).Once()
		userRepo.EXPECT().GetByID(ctx, receiver.ID).Return(receiver, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		categoryRepo.EXPECT().GetByID(ctx, catgry.ID).Return(catgry, nil).Once()
		expenseRepo.EXPECT().GetNextID().Return(expense.ID{Value: 1}).Once()
		expenseRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		p := usecase.CreateExpenseParams{
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   "transfer",
			Name:        "name",
			Amount:      100,
			Description: "description",
		}

		expns, err := createExpense(ctx, p)
		assert.Equal(t, expense.ID{Value: 1}, expns.ID)
		assert.Equal(t, "name", expns.Name)
		assert.Equal(t, 100, expns.Amount)
		assert.Equal(t, "description", expns.Description)
		assert.Equal(t, grp.ID, expns.GroupID)
		assert.Equal(t, catgry.ID, expns.CategoryID)
		assert.Equal(t, expense.SplitRatio{
			Payer:    0,
			Receiver: 100,
		}, expns.SplitRatio)
		assert.Equal(t, payer.ID, expns.PayerID)
		assert.Equal(t, receiver.ID, expns.ReceiverID)
		assert.Nil(t, err)
	})
}
