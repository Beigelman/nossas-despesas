package usecase

import (
	"context"
	"errors"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	mockrepository "github.com/Beigelman/ludaapi/internal/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAddUserToGroup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userRepo := mockrepository.NewMockUserRepository(t)
	groupRepo := mockrepository.NewMockGroupRepository(t)

	addUserToGroup := NewAddUserToGroup(userRepo, groupRepo)
	groupID := entity.GroupID{Value: 1}
	userID := entity.UserID{Value: 1}
	input := AddUserToGroupInput{
		GroupID: entity.GroupID{Value: 1},
		UserID:  entity.UserID{Value: 1},
	}
	group := entity.NewGroup(entity.GroupParams{ID: groupID})
	user := entity.NewUser(entity.UserParams{ID: userID})

	t.Run("successful scenario", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, input.UserID).Return(user, nil).Once()
		groupRepo.EXPECT().GetByID(ctx, input.GroupID).Return(group, nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		user, err := addUserToGroup(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, &input.GroupID, user.GroupID)
	})

	t.Run("group does not exist", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, input.GroupID).Return(nil, nil).Once()

		user, err := addUserToGroup(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group not found")
		assert.Nil(t, user)
	})

	t.Run("user does not exist", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, input.GroupID).Return(group, nil).Once()
		userRepo.EXPECT().GetByID(ctx, input.UserID).Return(nil, nil).Once()

		user, err := addUserToGroup(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		assert.Nil(t, user)
	})

	t.Run("error fetching group", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, input.GroupID).Return(nil, errors.New("db error")).Once()

		user, err := addUserToGroup(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "groupRepository.GetByID")
		assert.Nil(t, user)
	})

	t.Run("error fetching user", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, input.GroupID).Return(group, nil).Once()
		userRepo.EXPECT().GetByID(ctx, input.UserID).Return(nil, errors.New("db error")).Once()

		user, err := addUserToGroup(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "userRepository.GetByID")
		assert.Nil(t, user)
	})

	t.Run("error storing user", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, input.GroupID).Return(group, nil).Once()
		userRepo.EXPECT().GetByID(ctx, input.UserID).Return(user, nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("db error")).Once()

		user, err := addUserToGroup(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "userRepository.Store")
		assert.Nil(t, user)
	})
}
