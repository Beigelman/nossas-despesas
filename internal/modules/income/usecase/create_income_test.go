package usecase

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateIncome(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	incomeRepo := mockrepository.NewMockIncomeRepository(t)
	userRepo := mockrepository.NewMockUserRepository(t)

	user := entity.NewUser(entity.UserParams{
		ID:    entity.UserID{Value: 1},
		Name:  "My test user",
		Email: "email",
	})

	useCase := NewCreateIncome(userRepo, incomeRepo)
	params := CreateIncomeParams{
		Type:      income.Types.Salary,
		Amount:    100,
		UserID:    entity.UserID{Value: 1},
		CreatedAt: nil,
	}

	t.Run("getUserByID returns error", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, user.ID).Return(nil, errors.New("test error")).Once()
		inc, err := useCase(ctx, params)
		assert.Errorf(t, err, "userRepo.GetByID: test error")
		assert.Nil(t, inc)
	})

	t.Run("user does not exist", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, user.ID).Return(nil, nil).Once()
		inc, err := useCase(ctx, params)
		assert.Errorf(t, err, "user not found")
		assert.Nil(t, inc)
	})

	t.Run("store returns error", func(t *testing.T) {
		incomeRepo.EXPECT().GetNextID().Return(income.ID{Value: 1}).Once()
		userRepo.EXPECT().GetByID(ctx, user.ID).Return(user, nil).Once()
		incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		inc, err := useCase(ctx, params)
		assert.Errorf(t, err, "incomeRepo.Store: test error")
		assert.Nil(t, inc)
	})

	t.Run("success", func(t *testing.T) {
		incomeRepo.EXPECT().GetNextID().Return(income.ID{Value: 1}).Once()
		userRepo.EXPECT().GetByID(ctx, user.ID).Return(user, nil).Once()
		incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		inc, err := useCase(ctx, params)
		assert.NoError(t, err)
		assert.NotNil(t, inc)
		assert.Equal(t, entity.UserID{Value: 1}, inc.UserID)
		assert.Equal(t, income.Types.Salary, inc.Type)
		assert.Equal(t, 100, inc.Amount)
	})

}
