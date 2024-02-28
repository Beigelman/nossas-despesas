package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type CreateCategoryInput struct {
	Name            string
	Icon            string
	CategoryGroupID entity.CategoryGroupID
}

type CreateCategory func(ctx context.Context, input CreateCategoryInput) (*entity.Category, error)

func NewCreateCategory(categoryRepo repository.CategoryRepository) CreateCategory {
	return func(ctx context.Context, input CreateCategoryInput) (*entity.Category, error) {
		alreadyExists, err := categoryRepo.GetByName(ctx, input.Name)
		if err != nil {
			return nil, fmt.Errorf("categoryRepo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, except.ConflictError("category already exists")
		}

		categoryID := categoryRepo.GetNextID()

		category := entity.NewCategory(entity.CategoryParams{
			ID:              categoryID,
			Name:            input.Name,
			Icon:            input.Icon,
			CategoryGroupID: input.CategoryGroupID,
		})

		if err := categoryRepo.Store(ctx, category); err != nil {
			return nil, fmt.Errorf("categoryRepo.Store: %w", err)
		}

		return category, nil
	}
}
