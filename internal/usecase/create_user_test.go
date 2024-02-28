package usecase

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := mockrepository.NewMockUserRepository(t)
	existUser := entity.NewUser(entity.UserParams{
		ID:    entity.UserID{Value: 1},
		Name:  "My test user",
		Email: "my@email.com",
	})

	useCase := NewCreateUser(repo)
	params := CreateUserParams{
		Name:           "New user",
		Email:          "my@email.com",
		ProfilePicture: nil,
	}

	t.Run("getByName returns error", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, existUser.Email).Return(nil, errors.New("test error")).Once()
		user, err := useCase(ctx, params)
		assert.Errorf(t, err, "repo.GetByEmail(): test error")
		assert.Nil(t, user)
	})

	t.Run("user already exists", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, mock.Anything).Return(existUser, nil).Once()
		user, err := useCase(ctx, params)
		assert.Errorf(t, err, "email already exists")
		assert.Nil(t, user)
	})

	t.Run("Store returns error", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, mock.Anything).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		repo.EXPECT().GetNextID().Return(entity.UserID{Value: 1}).Once()
		user, err := useCase(ctx, params)
		assert.Errorf(t, err, "repo.Store: test error")
		assert.Nil(t, user)
	})

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, mock.Anything).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		repo.EXPECT().GetNextID().Return(entity.UserID{Value: 1}).Once()
		user, err := useCase(ctx, params)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, entity.UserID{Value: 1}, user.ID)
		assert.Equal(t, "New user", user.Name)
		assert.Equal(t, "my@email.com", user.Email)
	})

}
