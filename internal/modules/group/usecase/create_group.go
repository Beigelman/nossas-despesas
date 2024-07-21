package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type CreateGroupInput struct {
	Name   string
	UserID entity.UserID
}

type CreateGroup func(ctx context.Context, params CreateGroupInput) (*group.Group, error)

func NewCreateGroup(userRepo repository.UserRepository, groupRepo group.Repository) CreateGroup {
	return func(ctx context.Context, params CreateGroupInput) (*group.Group, error) {
		user, err := userRepo.GetByID(ctx, params.UserID)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if user.GroupID != nil {
			return nil, except.UnprocessableEntityError("user already in a group")
		}

		newGroup := group.NewGroup(group.Attributes{
			ID:   groupRepo.GetNextID(),
			Name: params.Name,
		})

		if err := groupRepo.Store(ctx, newGroup); err != nil {
			return nil, fmt.Errorf("groupRepo.Store: %w", err)
		}

		user.AssignGroup(newGroup.ID)

		if err := userRepo.Store(ctx, user); err != nil {
			return nil, fmt.Errorf("userRepo.Store: %w", err)
		}

		return newGroup, nil
	}
}
