package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/group/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInviteUserToGroup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mocks.NewMockuserRepository(t)
	groupRepo := mocks.NewMockgroupRepository(t)
	groupInviteRepo := mocks.NewMockgroupInviteRepository(t)
	emailProvider := mocks.NewMockserviceEmailProvider(t)

	grp := group.New(group.Attributes{
		ID:   group.ID{Value: 1},
		Name: "name",
	})

	invitee := user.New(user.Attributes{
		ID:      user.ID{Value: 1},
		Name:    "john",
		Email:   "john@email.com",
		GroupID: &grp.ID,
	})
	var invites []group.Invite
	for i := 0; i < 4; i++ {
		invites = append(invites, *group.NewInvite(group.InviteAttributes{
			ID:        group.InviteID{Value: i},
			GroupID:   grp.ID,
			Token:     "testToken",
			Email:     "john@email.com",
			ExpiresAt: time.Now(),
		}))
	}
	inviteUserToGroup := usecase.NewInviteUserToGroup(userRepo, groupRepo, groupInviteRepo, emailProvider)

	t.Run("should return error if groupRepo fails", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(nil, errors.New("test error")).Once()

		invite, err := inviteUserToGroup(ctx, usecase.InviteUserToGroupInput{
			GroupID: grp.ID,
			Email:   "john@email.com",
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "groupRepo.GetByID: test error")
	})

	t.Run("should return error if group not found", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(nil, nil).Once()

		invite, err := inviteUserToGroup(ctx, usecase.InviteUserToGroupInput{
			GroupID: grp.ID,
			Email:   "john@email.com",
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "group not found")
	})

	t.Run("should return error if userRepo fails", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, errors.New("test error")).Once()

		invite, err := inviteUserToGroup(ctx, usecase.InviteUserToGroupInput{
			GroupID: grp.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "userRepo.GetByEmail: test error")
	})

	t.Run("should return error if invitee has group already", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(invitee, nil).Once()

		invite, err := inviteUserToGroup(ctx, usecase.InviteUserToGroupInput{
			GroupID: grp.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "user already in a group")
	})

	t.Run("should return error if GetGroupInvitesByEmail fails", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetGroupInvitesByEmail(ctx, grp.ID, invitee.Email).Return(nil, errors.New("test error")).Once()

		invite, err := inviteUserToGroup(ctx, usecase.InviteUserToGroupInput{
			GroupID: grp.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "groupInviteRepo.GetByEmail: test error")
	})

	t.Run("should return error if GetGroupInvitesByEmail returns more than 3 invites in the last 48 hours", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetGroupInvitesByEmail(ctx, grp.ID, invitee.Email).Return(invites, nil).Once()

		invite, err := inviteUserToGroup(ctx, usecase.InviteUserToGroupInput{
			GroupID: grp.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "too many invites sent to this email recently")
	})

	t.Run("should return error if groupInviteRepo fails", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetGroupInvitesByEmail(ctx, grp.ID, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetNextID().Return(group.InviteID{Value: 1}).Once()
		groupInviteRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()

		invite, err := inviteUserToGroup(ctx, usecase.InviteUserToGroupInput{
			GroupID: grp.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, invite)
		assert.EqualError(t, err, "groupInviteRepo.Store: test error")
	})

	t.Run("happy path", func(t *testing.T) {
		groupRepo.EXPECT().GetByID(ctx, grp.ID).Return(grp, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetGroupInvitesByEmail(ctx, grp.ID, invitee.Email).Return(nil, nil).Once()
		groupInviteRepo.EXPECT().GetNextID().Return(group.InviteID{Value: 1}).Once()
		groupInviteRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()

		invite, err := inviteUserToGroup(ctx, usecase.InviteUserToGroupInput{
			GroupID: grp.ID,
			Email:   invitee.Email,
		})
		assert.Nil(t, err)
		assert.Equal(t, invitee.Email, invite.Email)
		assert.Equal(t, grp.ID, invite.GroupID)
	})
}
