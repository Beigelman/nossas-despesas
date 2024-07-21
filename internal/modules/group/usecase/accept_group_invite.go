package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	AcceptGroupInviteInput struct {
		Email string
		Token string
	}

	AcceptGroupInvite func(ctx context.Context, input AcceptGroupInviteInput) error
)

func NewAcceptGroupInvite(
	userRepository repository.UserRepository,
	groupInviteRepository group.InviteRepository,
) AcceptGroupInvite {
	return func(ctx context.Context, input AcceptGroupInviteInput) error {
		user, err := userRepository.GetByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("userRepository.GetByEmail: %w", err)
		}

		if user.GroupID != nil {
			return except.UnprocessableEntityError("user already in a group")
		}

		groupInvite, err := groupInviteRepository.GetByToken(ctx, input.Token)
		if err != nil {
			return fmt.Errorf("groupInviteRepository.GetByToken: %w", err)
		}

		if groupInvite == nil {
			return except.NotFoundError("invite not found")
		}

		if err := groupInvite.CheckStatus(); err != nil {
			return except.UnprocessableEntityError("invalid invite").SetInternal(err)
		}

		if groupInvite.Email != user.Email {
			return except.UnprocessableEntityError("invalid invite email")
		}

		user.AssignGroup(groupInvite.GroupID)

		if err := userRepository.Store(ctx, user); err != nil {
			return fmt.Errorf("userRepository.Store: %w", err)
		}

		return nil
	}
}
