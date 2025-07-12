package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/pkg/pubsub"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateExpensesFromScheduled(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	scheduledExpenseRepo := mocks.NewMockexpenseScheduledExpenseRepository(t)
	publisher := mocks.NewMockpubsubPublisher(t)

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

	lastGeneratedAt := civil.DateOf(time.Now().AddDate(0, 0, -31)) // 31 dias atrás

	scheduledExpense, err := expense.NewScheduledExpense(expense.ScheduledExpenseAttributes{
		ID:              expense.ScheduledExpenseID{Value: 1},
		Name:            "test scheduled expense",
		Amount:          100,
		Description:     "test description",
		GroupID:         grp.ID,
		CategoryID:      catgry.ID,
		SplitType:       expense.SplitTypes.Equal,
		PayerID:         payer.ID,
		ReceiverID:      receiver.ID,
		FrequencyInDays: 30,
		LastGeneratedAt: &lastGeneratedAt,
	})
	assert.NoError(t, err)

	generateExpensesFromScheduled := usecase.NewGenerateExpensesFromScheduledUseCase(scheduledExpenseRepo, publisher)

	t.Run("should return error if GetActiveScheduledExpenses fails", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetActiveScheduledExpenses(ctx).Return(nil, errors.New("database error")).Once()

		expensesCreated, err := generateExpensesFromScheduled(ctx)
		assert.Equal(t, 0, expensesCreated)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get active scheduled expenses")
	})

	t.Run("should return 0 when no scheduled expenses should generate", func(t *testing.T) {
		// Criar uma despesa agendada que não deve gerar (última geração recente)
		recentLastGeneratedAt := civil.DateOf(time.Now().AddDate(0, 0, -10)) // 10 dias atrás
		recentScheduledExpense, err := expense.NewScheduledExpense(expense.ScheduledExpenseAttributes{
			ID:              expense.ScheduledExpenseID{Value: 2},
			Name:            "recent scheduled expense",
			Amount:          100,
			Description:     "test description",
			GroupID:         grp.ID,
			CategoryID:      catgry.ID,
			SplitType:       expense.SplitTypes.Equal,
			PayerID:         payer.ID,
			ReceiverID:      receiver.ID,
			FrequencyInDays: 30,
			LastGeneratedAt: &recentLastGeneratedAt,
		})
		assert.NoError(t, err)

		scheduledExpenseRepo.EXPECT().GetActiveScheduledExpenses(ctx).Return([]expense.ScheduledExpense{*recentScheduledExpense}, nil).Once()

		expensesCreated, err := generateExpensesFromScheduled(ctx)
		assert.Equal(t, 0, expensesCreated)
		assert.NoError(t, err)
	})

	t.Run("should return error if BulkStore fails", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetActiveScheduledExpenses(ctx).Return([]expense.ScheduledExpense{*scheduledExpense}, nil).Once()
		scheduledExpenseRepo.EXPECT().BulkStore(ctx, mock.AnythingOfType("[]expense.ScheduledExpense")).Return(errors.New("bulk store error")).Once()

		expensesCreated, err := generateExpensesFromScheduled(ctx)
		assert.Equal(t, 0, expensesCreated)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to store scheduled expenses")
	})

	t.Run("should return error if Publish fails", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetActiveScheduledExpenses(ctx).Return([]expense.ScheduledExpense{*scheduledExpense}, nil).Once()
		scheduledExpenseRepo.EXPECT().BulkStore(ctx, mock.AnythingOfType("[]expense.ScheduledExpense")).Return(nil).Once()
		publisher.EXPECT().Publish(ctx, pubsub.ExpensesTopic, mock.AnythingOfType("pubsub.ExpenseEvent")).Return(errors.New("publish error")).Once()

		expensesCreated, err := generateExpensesFromScheduled(ctx)
		assert.Equal(t, 0, expensesCreated)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to publish expense created event")
	})

	t.Run("should generate expenses successfully", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetActiveScheduledExpenses(ctx).Return([]expense.ScheduledExpense{*scheduledExpense}, nil).Once()
		scheduledExpenseRepo.EXPECT().BulkStore(ctx, mock.AnythingOfType("[]expense.ScheduledExpense")).Return(nil).Once()
		publisher.EXPECT().Publish(ctx, pubsub.ExpensesTopic, mock.AnythingOfType("pubsub.ExpenseEvent")).Return(nil).Once()

		expensesCreated, err := generateExpensesFromScheduled(ctx)
		assert.Equal(t, 1, expensesCreated)
		assert.NoError(t, err)
	})

	t.Run("should generate multiple expenses successfully", func(t *testing.T) {
		// Criar segunda despesa agendada
		scheduledExpense2, err := expense.NewScheduledExpense(expense.ScheduledExpenseAttributes{
			ID:              expense.ScheduledExpenseID{Value: 3},
			Name:            "test scheduled expense 2",
			Amount:          200,
			Description:     "test description 2",
			GroupID:         grp.ID,
			CategoryID:      catgry.ID,
			SplitType:       expense.SplitTypes.Equal,
			PayerID:         payer.ID,
			ReceiverID:      receiver.ID,
			FrequencyInDays: 30,
			LastGeneratedAt: &lastGeneratedAt,
		})
		assert.NoError(t, err)

		scheduledExpenseRepo.EXPECT().GetActiveScheduledExpenses(ctx).Return([]expense.ScheduledExpense{*scheduledExpense, *scheduledExpense2}, nil).Once()
		scheduledExpenseRepo.EXPECT().BulkStore(ctx, mock.AnythingOfType("[]expense.ScheduledExpense")).Return(nil).Once()
		publisher.EXPECT().Publish(ctx, pubsub.ExpensesTopic, mock.AnythingOfType("pubsub.ExpenseEvent")).Return(nil).Times(2)

		expensesCreated, err := generateExpensesFromScheduled(ctx)
		assert.Equal(t, 2, expensesCreated)
		assert.NoError(t, err)
	})

	t.Run("should return 0 when no expenses to generate", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetActiveScheduledExpenses(ctx).Return([]expense.ScheduledExpense{}, nil).Once()

		expensesCreated, err := generateExpensesFromScheduled(ctx)
		assert.Equal(t, 0, expensesCreated)
		assert.NoError(t, err)
	})
}
