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

func TestCreateGroup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	groupRepo := mockrepository.NewMockGroupRepository(t)
	userRepo := mockrepository.NewMockUserRepository(t)
	user := entity.NewUser(entity.UserParams{
		ID:    entity.UserID{Value: 1},
		Name:  "test",
		Email: "",
	})

	useCase := NewCreateGroup(userRepo, groupRepo)

	t.Run("userRepo.GetByID returns error", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, user.ID).Return(nil, errors.New("test error")).Once()
		group, err := useCase(ctx, CreateGroupInput{
			Name:   "my new group",
			UserID: user.ID,
		})
		assert.Nil(t, group)
		assert.Errorf(t, err, "userRepo.GetByID: test error")
	})

	t.Run("user already in a group", func(t *testing.T) {
		user.GroupID = &entity.GroupID{Value: 1}
		userRepo.EXPECT().GetByID(ctx, user.ID).Return(user, nil).Once()
		group, err := useCase(ctx, CreateGroupInput{
			Name:   "my new group",
			UserID: user.ID,
		})
		assert.Nil(t, group)
		assert.Errorf(t, err, "user already in a group")
	})

	t.Run("groupRepo.Store returns error", func(t *testing.T) {
		user.GroupID = nil
		userRepo.EXPECT().GetByID(ctx, user.ID).Return(user, nil).Once()
		groupRepo.EXPECT().GetNextID().Return(entity.GroupID{Value: 1}).Once()
		groupRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		group, err := useCase(ctx, CreateGroupInput{
			Name:   "my new group",
			UserID: user.ID,
		})
		assert.Errorf(t, err, "groupRepo.Store: test error")
		assert.Nil(t, group)
	})

	t.Run("userRepo.Store returns error", func(t *testing.T) {
		user.GroupID = nil
		userRepo.EXPECT().GetByID(ctx, user.ID).Return(user, nil).Once()
		groupRepo.EXPECT().GetNextID().Return(entity.GroupID{Value: 1}).Once()
		groupRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		group, err := useCase(ctx, CreateGroupInput{
			Name:   "my new group",
			UserID: user.ID,
		})
		assert.Errorf(t, err, "userRepo.Store: test error")
		assert.Nil(t, group)
	})

	t.Run("success", func(t *testing.T) {
		user.GroupID = nil
		userRepo.EXPECT().GetByID(ctx, user.ID).Return(user, nil).Once()
		groupRepo.EXPECT().GetNextID().Return(entity.GroupID{Value: 1}).Once()
		groupRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		group, err := useCase(ctx, CreateGroupInput{
			Name:   "my new group",
			UserID: user.ID,
		})
		assert.NoError(t, err)
		assert.Equal(t, entity.GroupID{Value: 1}, group.ID)
		assert.Equal(t, "my new group", group.Name)
	})

}
