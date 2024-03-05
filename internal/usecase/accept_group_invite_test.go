package usecase

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestAcceptGroupInvite(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userRepo := mockrepository.NewMockUserRepository(t)
	groupInviteRepo := mockrepository.NewMockGroupInviteRepository(t)

	acceptGroupInvite := NewAcceptGroupInvite(userRepo, groupInviteRepo)
	groupID := entity.GroupID{Value: 1}
	userID := entity.UserID{Value: 1}
	input := AcceptGroupInviteInput{
		Email: "test@gmail.com",
		Token: "token",
	}
	groupInvite := entity.NewGroupInvite(entity.GroupInviteParams{
		ID:        entity.GroupInviteID{Value: 1},
		GroupID:   groupID,
		Token:     input.Token,
		Email:     input.Email,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	group := entity.NewGroup(entity.GroupParams{ID: groupID})
	user := entity.NewUser(entity.UserParams{ID: userID, GroupID: &group.ID})
	userWithOutGroup := entity.NewUser(entity.UserParams{ID: userID, Email: input.Email})

	t.Run("if user repo fails it returns error", func(t *testing.T) {
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(nil, errors.New("test error")).Once()
		assert.Error(t, acceptGroupInvite(ctx, input), "userRepository.GetByEmail: test error")
	})

	t.Run("user already have a group returns error", func(t *testing.T) {
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(user, nil).Once()
		assert.Error(t, acceptGroupInvite(ctx, input), "user already in a group")
	})

	t.Run("if groupInvite repo fails return error", func(t *testing.T) {
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(userWithOutGroup, nil).Once()
		groupInviteRepo.EXPECT().GetByToken(ctx, input.Token).Return(nil, errors.New("test error")).Once()
		assert.Error(t, acceptGroupInvite(ctx, input), "groupInviteRepository.GetByToken: test error")
	})

	t.Run("if groupInvite is nil return error", func(t *testing.T) {
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(userWithOutGroup, nil).Once()
		groupInviteRepo.EXPECT().GetByToken(ctx, input.Token).Return(nil, nil).Once()
		assert.Error(t, acceptGroupInvite(ctx, input), "invite not found")
	})

	t.Run("if groupInvite is invalid return error", func(t *testing.T) {
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(userWithOutGroup, nil).Once()
		groupInviteRepo.EXPECT().GetByToken(ctx, input.Token).Return(groupInvite, nil).Once()
		assert.Error(t, acceptGroupInvite(ctx, input), "invalid invite")
	})

	t.Run("if groupInvite email is invalid return error", func(t *testing.T) {
		assert.NoError(t, groupInvite.Sent())
		groupInvite.Email = "wrong@email.com"
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(userWithOutGroup, nil).Once()
		groupInviteRepo.EXPECT().GetByToken(ctx, input.Token).Return(groupInvite, nil).Once()
		assert.Error(t, acceptGroupInvite(ctx, input), "invalid invite email")
	})

	t.Run("if user repo fails to store return error", func(t *testing.T) {
		assert.NoError(t, groupInvite.Sent())
		groupInvite.Email = input.Email
		groupInviteRepo.EXPECT().GetByToken(ctx, input.Token).Return(groupInvite, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(userWithOutGroup, nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		assert.Error(t, acceptGroupInvite(ctx, input), "userRepository.Store: test error")
	})

	t.Run("if everything is ok return nil", func(t *testing.T) {
		assert.NoError(t, groupInvite.Sent())
		groupInvite.Email = input.Email
		userWithOutGroup.GroupID = nil
		groupInviteRepo.EXPECT().GetByToken(ctx, input.Token).Return(groupInvite, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(userWithOutGroup, nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		assert.NoError(t, acceptGroupInvite(ctx, input))
	})
}
