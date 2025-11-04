package usecase

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type CreateGroupInput struct {
	Name   string
	UserID user.ID
}

type CreateGroup func(ctx context.Context, params CreateGroupInput) (*group.Group, error)

func NewCreateGroup(userRepo user.Repository, groupRepo group.Repository) CreateGroup {
	return func(ctx context.Context, params CreateGroupInput) (*group.Group, error) {
		usr, err := userRepo.GetByID(ctx, params.UserID)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if usr.GroupID != nil {
			return nil, except.UnprocessableEntityError("user already in a group")
		}

		newGroup := group.New(group.Attributes{
			ID:   groupRepo.GetNextID(),
			Name: params.Name,
		})

		if err := groupRepo.Store(ctx, newGroup); err != nil {
			return nil, fmt.Errorf("groupRepo.Store: %w", err)
		}

		usr.AssignGroup(newGroup.ID)

		if err := userRepo.Store(ctx, usr); err != nil {
			return nil, fmt.Errorf("userRepo.Store: %w", err)
		}

		return newGroup, nil
	}
}
