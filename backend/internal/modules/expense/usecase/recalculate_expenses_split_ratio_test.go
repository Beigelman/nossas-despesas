package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRecalculateExpensesSplitRatio(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
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

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	proportionalExpense, err := expense.New(expense.Attributes{
		ID:          expense.ID{Value: 1},
		Name:        "test expense",
		Amount:      100,
		Description: "test description",
		GroupID:     grp.ID,
		CategoryID:  catgry.ID,
		SplitType:   expense.SplitTypes.Proportional,
		PayerID:     payer.ID,
		ReceiverID:  receiver.ID,
		SplitRatio:  expense.NewProportionalSplitRatio(6000, 4000), // 60% payer, 40% receiver
		CreatedAt:   &date,
	})
	assert.NoError(t, err)

	equalExpense, err := expense.New(expense.Attributes{
		ID:          expense.ID{Value: 2},
		Name:        "equal expense",
		Amount:      100,
		Description: "test description",
		GroupID:     grp.ID,
		CategoryID:  catgry.ID,
		SplitType:   expense.SplitTypes.Equal,
		PayerID:     payer.ID,
		ReceiverID:  receiver.ID,
		SplitRatio:  expense.NewEqualSplitRatio(),
		CreatedAt:   &date,
	})
	assert.NoError(t, err)

	payerIncome := income.New(income.Attributes{
		ID:        income.ID{Value: 1},
		UserID:    payer.ID,
		Amount:    6000,
		Type:      income.Types.Salary,
		CreatedAt: &date,
	})

	receiverIncome := income.New(income.Attributes{
		ID:        income.ID{Value: 2},
		UserID:    receiver.ID,
		Amount:    4000,
		Type:      income.Types.Salary,
		CreatedAt: &date,
	})

	recalculateExpensesSplitRatio := usecase.NewRecalculateExpensesSplitRatio(expenseRepo, incomeRepo)

	t.Run("should return error if GetByGroupDate fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByGroupDate(ctx, grp.ID, date).Return(nil, errors.New("database error")).Once()

		input := usecase.RecalculateExpensesSplitRatioInput{
			EventName: "test_event",
			GroupID:   grp.ID,
			Date:      date,
		}

		err := recalculateExpensesSplitRatio(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expensesRepo.GetByGroupDate")
	})

	t.Run("should return nil when no proportional expenses found", func(t *testing.T) {
		expenseRepo.EXPECT().GetByGroupDate(ctx, grp.ID, date).Return([]expense.Expense{*equalExpense}, nil).Once()

		input := usecase.RecalculateExpensesSplitRatioInput{
			EventName: "test_event",
			GroupID:   grp.ID,
			Date:      date,
		}

		err := recalculateExpensesSplitRatio(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("should return error if GetUserMonthlyIncomes fails for payer", func(t *testing.T) {
		expenseRepo.EXPECT().GetByGroupDate(ctx, grp.ID, date).Return([]expense.Expense{*proportionalExpense}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, &date).Return(nil, errors.New("income error")).Once()

		input := usecase.RecalculateExpensesSplitRatioInput{
			EventName: "test_event",
			GroupID:   grp.ID,
			Date:      date,
		}

		err := recalculateExpensesSplitRatio(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no incomes found for user")
	})

	t.Run("should return error if GetUserMonthlyIncomes returns nil for payer", func(t *testing.T) {
		expenseRepo.EXPECT().GetByGroupDate(ctx, grp.ID, date).Return([]expense.Expense{*proportionalExpense}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, &date).Return(nil, nil).Once()

		input := usecase.RecalculateExpensesSplitRatioInput{
			EventName: "test_event",
			GroupID:   grp.ID,
			Date:      date,
		}

		err := recalculateExpensesSplitRatio(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no incomes found for user")
	})

	t.Run("should return error if GetUserMonthlyIncomes fails for receiver", func(t *testing.T) {
		expenseRepo.EXPECT().GetByGroupDate(ctx, grp.ID, date).Return([]expense.Expense{*proportionalExpense}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, &date).Return([]income.Income{*payerIncome}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, receiver.ID, &date).Return(nil, errors.New("income error")).Once()

		input := usecase.RecalculateExpensesSplitRatioInput{
			EventName: "test_event",
			GroupID:   grp.ID,
			Date:      date,
		}

		err := recalculateExpensesSplitRatio(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no incomes found for user")
	})

	t.Run("should return error if BulkStore fails", func(t *testing.T) {
		expenseRepo.EXPECT().GetByGroupDate(ctx, grp.ID, date).Return([]expense.Expense{*proportionalExpense}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, &date).Return([]income.Income{*payerIncome}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, receiver.ID, &date).Return([]income.Income{*receiverIncome}, nil).Once()
		expenseRepo.EXPECT().BulkStore(ctx, mock.AnythingOfType("[]expense.Expense")).Return(errors.New("bulk store error")).Once()

		input := usecase.RecalculateExpensesSplitRatioInput{
			EventName: "test_event",
			GroupID:   grp.ID,
			Date:      date,
		}

		err := recalculateExpensesSplitRatio(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expense.BulkStore")
	})

	t.Run("should recalculate expenses split ratio successfully", func(t *testing.T) {
		expenseRepo.EXPECT().GetByGroupDate(ctx, grp.ID, date).Return([]expense.Expense{*proportionalExpense}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, &date).Return([]income.Income{*payerIncome}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, receiver.ID, &date).Return([]income.Income{*receiverIncome}, nil).Once()
		expenseRepo.EXPECT().BulkStore(ctx, mock.AnythingOfType("[]expense.Expense")).Return(nil).Once()

		input := usecase.RecalculateExpensesSplitRatioInput{
			EventName: "test_event",
			GroupID:   grp.ID,
			Date:      date,
		}

		err := recalculateExpensesSplitRatio(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("should recalculate multiple proportional expenses successfully", func(t *testing.T) {
		// Criar segunda despesa proporcional
		proportionalExpense2, err := expense.New(expense.Attributes{
			ID:          expense.ID{Value: 3},
			Name:        "test expense 2",
			Amount:      200,
			Description: "test description 2",
			GroupID:     grp.ID,
			CategoryID:  catgry.ID,
			SplitType:   expense.SplitTypes.Proportional,
			PayerID:     payer.ID,
			ReceiverID:  receiver.ID,
			SplitRatio:  expense.NewProportionalSplitRatio(6000, 4000),
			CreatedAt:   &date,
		})
		assert.NoError(t, err)

		expenseRepo.EXPECT().GetByGroupDate(ctx, grp.ID, date).Return([]expense.Expense{*proportionalExpense, *proportionalExpense2}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, &date).Return([]income.Income{*payerIncome}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, receiver.ID, &date).Return([]income.Income{*receiverIncome}, nil).Once()
		expenseRepo.EXPECT().BulkStore(ctx, mock.AnythingOfType("[]expense.Expense")).Return(nil).Once()

		input := usecase.RecalculateExpensesSplitRatioInput{
			EventName: "test_event",
			GroupID:   grp.ID,
			Date:      date,
		}

		err = recalculateExpensesSplitRatio(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("should handle mixed expense types correctly", func(t *testing.T) {
		// Misturar despesas proporcionais e iguais
		expenseRepo.EXPECT().GetByGroupDate(ctx, grp.ID, date).Return([]expense.Expense{*proportionalExpense, *equalExpense}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, payer.ID, &date).Return([]income.Income{*payerIncome}, nil).Once()
		incomeRepo.EXPECT().GetUserMonthlyIncomes(ctx, receiver.ID, &date).Return([]income.Income{*receiverIncome}, nil).Once()
		expenseRepo.EXPECT().BulkStore(ctx, mock.AnythingOfType("[]expense.Expense")).Return(nil).Once()

		input := usecase.RecalculateExpensesSplitRatioInput{
			EventName: "test_event",
			GroupID:   grp.ID,
			Date:      date,
		}

		err := recalculateExpensesSplitRatio(ctx, input)
		assert.NoError(t, err)
	})
}
