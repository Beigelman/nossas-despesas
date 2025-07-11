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

func TestCreateIncome(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mocks.NewMockuserRepository(t)
	incomeRepo := mocks.NewMockincomeRepository(t)
	publisher := mocks.NewMockpubsubPublisher(t)

	usr := user.New(user.Attributes{
		ID:    user.ID{Value: 1},
		Name:  "My test user",
		Email: "email",
	})

	params := usecase.CreateIncomeParams{
		Type:      income.Types.Salary,
		Amount:    100,
		UserID:    user.ID{Value: 1},
		GroupID:   group.ID{Value: 1},
		CreatedAt: nil,
	}

	useCase := usecase.NewCreateIncome(userRepo, incomeRepo, publisher)

	t.Run("getUserByID returns error", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(nil, errors.New("test error")).Once()
		inc, err := useCase(ctx, params)
		assert.ErrorContains(t, err, "userRepo.GetByID: test error")
		assert.Nil(t, inc)
	})

	t.Run("user does not exist", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(nil, nil).Once()
		inc, err := useCase(ctx, params)
		assert.ErrorContains(t, err, "user not found")
		assert.Nil(t, inc)
	})

	t.Run("store returns error", func(t *testing.T) {
		incomeRepo.EXPECT().GetNextID().Return(income.ID{Value: 1}).Once()
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(usr, nil).Once()
		incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		inc, err := useCase(ctx, params)
		assert.ErrorContains(t, err, "incomeRepo.Store: test error")
		assert.Nil(t, inc)
	})

	t.Run("success", func(t *testing.T) {
		incomeRepo.EXPECT().GetNextID().Return(income.ID{Value: 1}).Once()
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(usr, nil).Once()
		incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		publisher.EXPECT().Publish(ctx, pubsub.IncomesTopic, mock.Anything).Return(nil).Once()
		inc, err := useCase(ctx, params)
		assert.NoError(t, err)
		assert.NotNil(t, inc)
		assert.Equal(t, user.ID{Value: 1}, inc.UserID)
		assert.Equal(t, income.Types.Salary, inc.Type)
		assert.Equal(t, 100, inc.Amount)
	})

	t.Run("publish returns error (should not fail usecase)", func(t *testing.T) {
		incomeRepo.EXPECT().GetNextID().Return(income.ID{Value: 1}).Once()
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(usr, nil).Once()
		incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		publisher.EXPECT().Publish(ctx, pubsub.IncomesTopic, mock.Anything).Return(errors.New("publish error")).Once()
		inc, err := useCase(ctx, params)
		assert.NoError(t, err)
		assert.NotNil(t, inc)
	})
}
