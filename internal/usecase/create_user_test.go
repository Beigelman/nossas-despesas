package usecase

import (
	"context"
	"errors"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := mocks.NewMockUserRepository(t)
	existUser := entity.NewUser(entity.UserParams{
		ID:    entity.UserID{Value: 1},
		Name:  "My test group",
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
		group, err := useCase(ctx, params)
		assert.Errorf(t, err, "repo.GetByEmail(): test error")
		assert.Nil(t, group)
	})

	t.Run("group already exists", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, mock.Anything).Return(existUser, nil).Once()
		group, err := useCase(ctx, params)
		assert.Errorf(t, err, "email already exists")
		assert.Nil(t, group)
	})

	t.Run("Store returns error", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, mock.Anything).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		repo.EXPECT().GetNextID().Return(entity.UserID{Value: 1}).Once()
		group, err := useCase(ctx, params)
		assert.Errorf(t, err, "repo.Store: test error")
		assert.Nil(t, group)
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
