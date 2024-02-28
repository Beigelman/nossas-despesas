package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type CreateCategoryGroupInput struct {
	Name string
	Icon string
}

type CreateCategoryGroup func(ctx context.Context, input CreateCategoryGroupInput) (*entity.CategoryGroup, error)

func NewCreateCategoryGroup(categoryGroupRepo repository.CategoryGroupRepository) CreateCategoryGroup {
	return func(ctx context.Context, input CreateCategoryGroupInput) (*entity.CategoryGroup, error) {
		alreadyExists, err := categoryGroupRepo.GetByName(ctx, input.Name)
		if err != nil {
			return nil, fmt.Errorf("categoryGroupRepo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, except.ConflictError("category already exists")
		}

		categoryGroupID := categoryGroupRepo.GetNextID()

		categoryGroup := entity.NewCategoryGroup(entity.CategoryGroupParams{
			ID:   categoryGroupID,
			Name: input.Name,
			Icon: input.Icon,
		})

		if err := categoryGroupRepo.Store(ctx, categoryGroup); err != nil {
			return nil, fmt.Errorf("categoryGroupRepo.Store: %w", err)
		}

		return categoryGroup, nil
	}
}
