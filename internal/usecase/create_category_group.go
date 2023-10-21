package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
)

type CreateCategoryGroupInput struct {
	Name string
	Icon string
}

type CreateCategoryGroup func(ctx context.Context, input CreateCategoryGroupInput) (*entity.CategoryGroup, error)

func NewCreateCategoryGroup(repo repository.CategoryGroupRepository) CreateCategoryGroup {
	return func(ctx context.Context, input CreateCategoryGroupInput) (*entity.CategoryGroup, error) {
		alreadyExists, err := repo.GetByName(ctx, input.Name)
		if err != nil {
			return nil, fmt.Errorf("repo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, except.ConflictError("category already exists")
		}

		categoryGroupID := repo.GetNextID()

		categoryGroup := entity.NewCategoryGroup(entity.CategoryGroupParams{
			ID:   categoryGroupID,
			Name: input.Name,
			Icon: input.Icon,
		})

		if err := repo.Store(ctx, categoryGroup); err != nil {
			return nil, fmt.Errorf("repo.Store: %w", err)
		}

		return categoryGroup, nil
	}
}
