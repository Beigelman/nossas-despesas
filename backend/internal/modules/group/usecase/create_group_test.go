package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/group/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
)

func TestCreateGroup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	groupRepo := mocks.NewMockgroupRepository(t)
	userRepo := mocks.NewMockuserRepository(t)
	usr := user.New(user.Attributes{
		ID:    user.ID{Value: 1},
		Name:  "test",
		Email: "",
	})

	useCase := usecase.NewCreateGroup(userRepo, groupRepo)

	t.Run("userRepo.GetByID returns error", func(t *testing.T) {
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(nil, errors.New("test error")).Once()
		grp, err := useCase(ctx, usecase.CreateGroupInput{
			Name:   "my new group",
			UserID: usr.ID,
		})
		assert.Nil(t, grp)
		assert.Errorf(t, err, "userRepo.GetByID: test error")
	})

	t.Run("user already in a group", func(t *testing.T) {
		usr.GroupID = &group.ID{Value: 1}
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(usr, nil).Once()
		grp, err := useCase(ctx, usecase.CreateGroupInput{
			Name:   "my new group",
			UserID: usr.ID,
		})
		assert.Nil(t, grp)
		assert.Errorf(t, err, "user already in a group")
	})

	t.Run("groupRepo.Store returns error", func(t *testing.T) {
		usr.GroupID = nil
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(usr, nil).Once()
		groupRepo.EXPECT().GetNextID().Return(group.ID{Value: 1}).Once()
		groupRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		grp, err := useCase(ctx, usecase.CreateGroupInput{
			Name:   "my new group",
			UserID: usr.ID,
		})
		assert.Errorf(t, err, "groupRepo.Store: test error")
		assert.Nil(t, grp)
	})

	t.Run("userRepo.Store returns error", func(t *testing.T) {
		usr.GroupID = nil
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(usr, nil).Once()
		groupRepo.EXPECT().GetNextID().Return(group.ID{Value: 1}).Once()
		groupRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		grp, err := useCase(ctx, usecase.CreateGroupInput{
			Name:   "my new group",
			UserID: usr.ID,
		})
		assert.Errorf(t, err, "userRepo.Store: test error")
		assert.Nil(t, grp)
	})

	t.Run("success", func(t *testing.T) {
		usr.GroupID = nil
		userRepo.EXPECT().GetByID(ctx, usr.ID).Return(usr, nil).Once()
		groupRepo.EXPECT().GetNextID().Return(group.ID{Value: 1}).Once()
		groupRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		grp, err := useCase(ctx, usecase.CreateGroupInput{
			Name:   "my new group",
			UserID: usr.ID,
		})
		assert.NoError(t, err)
		assert.Equal(t, group.ID{Value: 1}, grp.ID)
		assert.Equal(t, "my new group", grp.Name)
	})
}
