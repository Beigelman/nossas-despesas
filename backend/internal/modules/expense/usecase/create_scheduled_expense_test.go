package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
)

func TestCreateScheduledExpense(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	scheduledExpenseRepo := mocks.NewMockexpenseScheduledExpenseRepository(t)

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

	lastGeneratedAt := civil.DateOf(time.Now())

	createScheduledExpense := usecase.NewCreateScheduledExpense(scheduledExpenseRepo)

	t.Run("should return error if Store fails", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetNextID().Return(expense.ScheduledExpenseID{Value: 1}).Once()
		scheduledExpenseRepo.EXPECT().Store(ctx, mock.AnythingOfType("*expense.ScheduledExpense")).Return(errors.New("store error")).Once()

		input := usecase.CreateScheduledExpenseInput{
			Name:            "test expense",
			Amount:          100,
			Description:     "test description",
			GroupID:         grp.ID,
			CategoryID:      catgry.ID,
			SplitType:       expense.SplitTypes.Equal,
			PayerID:         payer.ID,
			ReceiverID:      receiver.ID,
			FrequencyInDays: 30,
			LastGeneratedAt: &lastGeneratedAt,
		}

		err := createScheduledExpense(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to store scheduled expense")
	})

	t.Run("should create scheduled expense successfully", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetNextID().Return(expense.ScheduledExpenseID{Value: 1}).Once()
		scheduledExpenseRepo.EXPECT().Store(ctx, mock.AnythingOfType("*expense.ScheduledExpense")).Return(nil).Once()

		input := usecase.CreateScheduledExpenseInput{
			Name:            "test expense",
			Amount:          100,
			Description:     "test description",
			GroupID:         grp.ID,
			CategoryID:      catgry.ID,
			SplitType:       expense.SplitTypes.Equal,
			PayerID:         payer.ID,
			ReceiverID:      receiver.ID,
			FrequencyInDays: 30,
			LastGeneratedAt: &lastGeneratedAt,
		}

		err := createScheduledExpense(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("should create scheduled expense with proportional split type", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetNextID().Return(expense.ScheduledExpenseID{Value: 1}).Once()
		scheduledExpenseRepo.EXPECT().Store(ctx, mock.AnythingOfType("*expense.ScheduledExpense")).Return(nil).Once()

		input := usecase.CreateScheduledExpenseInput{
			Name:            "test expense",
			Amount:          100,
			Description:     "test description",
			GroupID:         grp.ID,
			CategoryID:      catgry.ID,
			SplitType:       expense.SplitTypes.Proportional,
			PayerID:         payer.ID,
			ReceiverID:      receiver.ID,
			FrequencyInDays: 30,
			LastGeneratedAt: &lastGeneratedAt,
		}

		err := createScheduledExpense(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("should create scheduled expense with transfer split type", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetNextID().Return(expense.ScheduledExpenseID{Value: 1}).Once()
		scheduledExpenseRepo.EXPECT().Store(ctx, mock.AnythingOfType("*expense.ScheduledExpense")).Return(nil).Once()

		input := usecase.CreateScheduledExpenseInput{
			Name:            "test expense",
			Amount:          100,
			Description:     "test description",
			GroupID:         grp.ID,
			CategoryID:      catgry.ID,
			SplitType:       expense.SplitTypes.Transfer,
			PayerID:         payer.ID,
			ReceiverID:      receiver.ID,
			FrequencyInDays: 30,
			LastGeneratedAt: &lastGeneratedAt,
		}

		err := createScheduledExpense(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("should return error if validation fails", func(t *testing.T) {
		scheduledExpenseRepo.EXPECT().GetNextID().Return(expense.ScheduledExpenseID{Value: 1}).Once()

		input := usecase.CreateScheduledExpenseInput{
			Name:            "", // Nome vazio deve causar erro de validação
			Amount:          100,
			Description:     "test description",
			GroupID:         grp.ID,
			CategoryID:      catgry.ID,
			SplitType:       expense.SplitTypes.Equal,
			PayerID:         payer.ID,
			ReceiverID:      receiver.ID,
			FrequencyInDays: 30,
			LastGeneratedAt: &lastGeneratedAt,
		}

		err := createScheduledExpense(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create scheduled expense")
	})
}
