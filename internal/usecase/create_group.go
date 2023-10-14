package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
)

type CreateGroup func(ctx context.Context, name string) (*entity.Group, error)

func NewCreateGroup(repo repository.GroupRepository) CreateGroup {
	return func(ctx context.Context, name string) (*entity.Group, error) {
		alreadyExists, err := repo.GetByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("repo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, fmt.Errorf("group already exists")
		}

		groupID := repo.GetNextID()

		group := entity.NewGroup(entity.GroupParams{
			ID:   groupID,
			Name: name,
		})

		if err := repo.Store(ctx, group); err != nil {
			return nil, fmt.Errorf("repo.Store: %w", err)
		}

		return group, nil
	}
}
