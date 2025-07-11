package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/pubsub"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteIncome(t *testing.T) {
	ctx := context.Background()
	incomeRepo := mocks.NewMockincomeRepository(t)
	userRepo := mocks.NewMockuserRepository(t)
	publisher := mocks.NewMockpubsubPublisher(t)

	usr := user.New(user.Attributes{
		ID:    user.ID{Value: 1},
		Name:  "Test User",
		Email: "email",
	})
	usr.AddFlag(user.EDIT_PARTNER_INCOME)

	inc := income.New(income.Attributes{
		ID:     income.ID{Value: 1},
		UserID: usr.ID,
		Amount: 100,
		Type:   income.Types.Salary,
	})

	params := usecase.DeleteIncomeParams{
		ID:      income.ID{Value: 1},
		UserID:  usr.ID,
		GroupID: group.ID{Value: 1},
	}

	useCase := usecase.NewDeleteIncome(incomeRepo, userRepo, publisher)

	t.Run("incomeRepo.GetByID returns error", func(t *testing.T) {
		incomeRepo.EXPECT().GetByID(ctx, params.ID).Return(nil, errors.New("test error")).Once()
		res, err := useCase(ctx, params)
		assert.ErrorContains(t, err, "incomeRepo.GetByID: test error")
		assert.Nil(t, res)
	})

	t.Run("userRepo.GetByID returns error", func(t *testing.T) {
		incomeRepo.EXPECT().GetByID(ctx, params.ID).Return(inc, nil).Once()
		userRepo.EXPECT().GetByID(ctx, params.UserID).Return(nil, errors.New("test error")).Once()
		res, err := useCase(ctx, params)
		assert.ErrorContains(t, err, "userRepo.GetByID: test error")
		assert.Nil(t, res)
	})

	t.Run("user not found", func(t *testing.T) {
		incomeRepo.EXPECT().GetByID(ctx, params.ID).Return(inc, nil).Once()
		userRepo.EXPECT().GetByID(ctx, params.UserID).Return(nil, nil).Once()
		res, err := useCase(ctx, params)
		assert.ErrorContains(t, err, "user not found")
		assert.Nil(t, res)
	})

	t.Run("user mismatch", func(t *testing.T) {
		incomeRepo.EXPECT().GetByID(ctx, params.ID).Return(inc, nil).Once()
		otherUser := user.New(user.Attributes{ID: user.ID{Value: 2}})
		userRepo.EXPECT().GetByID(ctx, mock.Anything).Return(otherUser, nil).Once()
		res, err := useCase(ctx, params)
		assert.ErrorContains(t, err, "user mismatch")
		assert.Nil(t, res)
	})

	t.Run("incomeRepo.Store returns error", func(t *testing.T) {
		incomeRepo.EXPECT().GetByID(ctx, params.ID).Return(inc, nil).Once()
		userRepo.EXPECT().GetByID(ctx, params.UserID).Return(usr, nil).Once()
		incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("store error")).Once()
		res, err := useCase(ctx, params)
		assert.ErrorContains(t, err, "incomeRepo.Store: store error")
		assert.Nil(t, res)
	})

	t.Run("success", func(t *testing.T) {
		incomeRepo.EXPECT().GetByID(ctx, params.ID).Return(inc, nil).Once()
		userRepo.EXPECT().GetByID(ctx, params.UserID).Return(usr, nil).Once()
		incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		publisher.EXPECT().Publish(ctx, pubsub.IncomesTopic, mock.Anything).Return(nil).Once()
		res, err := useCase(ctx, params)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("publish returns error (should not fail usecase)", func(t *testing.T) {
		incomeRepo.EXPECT().GetByID(ctx, params.ID).Return(inc, nil).Once()
		userRepo.EXPECT().GetByID(ctx, params.UserID).Return(usr, nil).Once()
		incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		publisher.EXPECT().Publish(ctx, pubsub.IncomesTopic, mock.Anything).Return(errors.New("publish error")).Once()
		res, err := useCase(ctx, params)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}
