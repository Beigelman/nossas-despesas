package usecase

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type CreateCategoryInput struct {
	Name            string
	Icon            string
	CategoryGroupID category.GroupID
}

type CreateCategory func(ctx context.Context, input CreateCategoryInput) (*category.Category, error)

func NewCreateCategory(categoryRepo category.Repository) CreateCategory {
	return func(ctx context.Context, input CreateCategoryInput) (*category.Category, error) {
		alreadyExists, err := categoryRepo.GetByName(ctx, input.Name)
		if err != nil {
			return nil, fmt.Errorf("categoryRepo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, except.ConflictError("category already exists")
		}

		categoryID := categoryRepo.GetNextID()

		newCategory := category.New(category.Attributes{
			ID:              categoryID,
			Name:            input.Name,
			Icon:            input.Icon,
			CategoryGroupID: input.CategoryGroupID,
		})

		if err := categoryRepo.Store(ctx, newCategory); err != nil {
			return nil, fmt.Errorf("categoryRepo.Store: %w", err)
		}

		return newCategory, nil
	}
}
