package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/modules/user/usecase"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := mocks.NewMockuserRepository(t)
	existUser := user.New(user.Attributes{
		ID:    user.ID{Value: 1},
		Name:  "My test user",
		Email: "my@email.com",
	})

	useCase := usecase.NewCreateUser(repo)
	params := usecase.CreateUserParams{
		Name:           "New user",
		Email:          "my@email.com",
		ProfilePicture: nil,
	}

	t.Run("getByName returns error", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, existUser.Email).Return(nil, errors.New("test error")).Once()
		usr, err := useCase(ctx, params)
		assert.Errorf(t, err, "repo.GetByEmail(): test error")
		assert.Nil(t, usr)
	})

	t.Run("user already exists", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, mock.Anything).Return(existUser, nil).Once()
		usr, err := useCase(ctx, params)
		assert.Errorf(t, err, "email already exists")
		assert.Nil(t, usr)
	})

	t.Run("Store returns error", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, mock.Anything).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		repo.EXPECT().GetNextID().Return(user.ID{Value: 1}).Once()
		usr, err := useCase(ctx, params)
		assert.Errorf(t, err, "repo.Store: test error")
		assert.Nil(t, usr)
	})

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetByEmail(ctx, mock.Anything).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		repo.EXPECT().GetNextID().Return(user.ID{Value: 1}).Once()
		usr, err := useCase(ctx, params)
		assert.NoError(t, err)
		assert.NotNil(t, usr)
		assert.Equal(t, user.ID{Value: 1}, usr.ID)
		assert.Equal(t, "New user", usr.Name)
		assert.Equal(t, "my@email.com", usr.Email)
	})
}
