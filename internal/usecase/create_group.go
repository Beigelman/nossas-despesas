package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type CreateGroupInput struct {
	Name   string
	UserID entity.UserID
}

type CreateGroup func(ctx context.Context, params CreateGroupInput) (*entity.Group, error)

func NewCreateGroup(userRepo repository.UserRepository, groupRepo repository.GroupRepository) CreateGroup {
	return func(ctx context.Context, params CreateGroupInput) (*entity.Group, error) {
		user, err := userRepo.GetByID(ctx, params.UserID)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if user.GroupID != nil {
			return nil, except.UnprocessableEntityError("user already in a group")
		}

		group := entity.NewGroup(entity.GroupParams{
			ID:   groupRepo.GetNextID(),
			Name: params.Name,
		})

		if err := groupRepo.Store(ctx, group); err != nil {
			return nil, fmt.Errorf("groupRepo.Store: %w", err)
		}

		user.AssignGroup(group.ID)

		if err := userRepo.Store(ctx, user); err != nil {
			return nil, fmt.Errorf("userRepo.Store: %w", err)
		}

		return group, nil
	}
}
