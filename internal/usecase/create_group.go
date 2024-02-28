package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type CreateGroup func(ctx context.Context, name string) (*entity.Group, error)

func NewCreateGroup(groupRepo repository.GroupRepository) CreateGroup {
	return func(ctx context.Context, name string) (*entity.Group, error) {
		alreadyExists, err := groupRepo.GetByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("groupRepo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, except.ConflictError("group already exists")
		}

		groupID := groupRepo.GetNextID()

		group := entity.NewGroup(entity.GroupParams{
			ID:   groupID,
			Name: name,
		})

		if err := groupRepo.Store(ctx, group); err != nil {
			return nil, fmt.Errorf("groupRepo.Store: %w", err)
		}

		return group, nil
	}
}
