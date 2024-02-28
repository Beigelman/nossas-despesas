package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	AddUserToGroupInput struct {
		GroupID entity.GroupID
		UserID  entity.UserID
	}

	AddUserToGroup func(ctx context.Context, input AddUserToGroupInput) (*entity.User, error)
)

func NewAddUserToGroup(userRepository repository.UserRepository, groupRepository repository.GroupRepository) AddUserToGroup {
	return func(ctx context.Context, input AddUserToGroupInput) (*entity.User, error) {
		group, err := groupRepository.GetByID(ctx, input.GroupID)
		if err != nil {
			return nil, fmt.Errorf("groupRepository.GetByID: %w", err)
		}

		if group == nil {
			return nil, except.NotFoundError("group not found")
		}

		user, err := userRepository.GetByID(ctx, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("userRepository.GetByID: %w", err)
		}

		if user == nil {
			return nil, except.NotFoundError("user not found")
		}

		user.AssignGroup(group.ID)

		if err := userRepository.Store(ctx, user); err != nil {
			return nil, fmt.Errorf("userRepository.Store: %w", err)
		}

		return user, nil
	}
}
