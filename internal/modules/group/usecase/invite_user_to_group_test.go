package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
	mockservice "github.com/Beigelman/nossas-despesas/internal/tests/mocks/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestInviteUserToGroup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mockrepository.NewMockUserRepository(t)
	groupRepo := mockrepository.NewMockGroupRepository(t)
	groupInviteRepo := mockrepository.NewMockGroupInviteRepository(t)
	emailProvider := mockservice.NewMockEmailProvider(t)

	group := group.NewGroup(group.Attributes{
		ID:   group.ID{Value: 1},
		Name: "name",
	})

	invitee := entity.NewUser(entity.UserParams{
		ID:      entity.UserID{Value: 1},
		Name:    "john",
		Email:   "john@email.com",
		GroupID: &group.ID,
	})
	var invites []group.GroupInvite
	for i := 0; i < 4; i++ {
		invites = append(invites, *group.NewGroupInvite(group.GroupInviteParams{
			ID:        group.GroupInviteID{Value: i},
			GroupID:   group.ID,
			Token:     "testToken",
			Email:     "john@email.com",
			ExpiresAt: time.Now(),
		}))
	}
	inviteUserToGroup := NewInviteUserToGroup(userRepo, groupRepo, groupInviteRepo, emailProvider)

	t.Run("should return error if groupRepo fails", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, group.ID).Return(nil, errors.New("test error")).Once()

		invite, err := inviteUserToGroup(ctx, InviteUserToGroupInput{
			GroupID: group.ID,
			Email:   "john@email.com",
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "groupRepo.GetByID: test error")
	})

	t.Run("should return error if group not found", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, group.ID).Return(nil, nil).Once()

		invite, err := inviteUserToGroup(ctx, InviteUserToGroupInput{
			GroupID: group.ID,
			Email:   "john@email.com",
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "group not found")
	})

	t.Run("should return error if userRepo fails", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, group.ID).Return(group, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, errors.New("test error")).Once()

		invite, err := inviteUserToGroup(ctx, InviteUserToGroupInput{
			GroupID: group.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "userRepo.GetByEmail: test error")
	})

	t.Run("should return error if invitee has group already", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, group.ID).Return(group, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(invitee, nil).Once()

		invite, err := inviteUserToGroup(ctx, InviteUserToGroupInput{
			GroupID: group.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "user already in a group")
	})

	t.Run("should return error if GetGroupInvitesByEmail fails", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, group.ID).Return(group, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetGroupInvitesByEmail(ctx, group.ID, invitee.Email).Return(nil, errors.New("test error")).Once()

		invite, err := inviteUserToGroup(ctx, InviteUserToGroupInput{
			GroupID: group.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "groupInviteRepo.GetByEmail: test error")
	})

	t.Run("should return error if GetGroupInvitesByEmail returns more than 3 invites in the last 48 hours", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, group.ID).Return(group, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetGroupInvitesByEmail(ctx, group.ID, invitee.Email).Return(invites, nil).Once()

		invite, err := inviteUserToGroup(ctx, InviteUserToGroupInput{
			GroupID: group.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "too many invites sent to this email recently")
	})

	t.Run("should return error if groupInviteRepo fails", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, group.ID).Return(group, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetGroupInvitesByEmail(ctx, group.ID, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetNextID().Return(group.GroupInviteID{Value: 1}).Once()
		groupInviteRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()

		invite, err := inviteUserToGroup(ctx, InviteUserToGroupInput{
			GroupID: group.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "groupInviteRepo.Store: test error")
	})

	t.Run("happy path", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, group.ID).Return(group, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetGroupInvitesByEmail(ctx, group.ID, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetNextID().Return(group.GroupInviteID{Value: 1}).Once()
		groupInviteRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		invite, err := inviteUserToGroup(ctx, InviteUserToGroupInput{
			GroupID: group.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, err)
		assert.Equal(t, invitee.Email, invite.Email)
		assert.Equal(t, group.ID, invite.GroupID)
	})
}
