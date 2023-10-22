package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
)

type CreateCategoryInput struct {
	Name            string
	Icon            string
	CategoryGroupID entity.CategoryGroupID
}

type CreateCategory func(ctx context.Context, input CreateCategoryInput) (*entity.Category, error)

func NewCreateCategory(repo repository.CategoryRepository) CreateCategory {
	return func(ctx context.Context, input CreateCategoryInput) (*entity.Category, error) {
		alreadyExists, err := repo.GetByName(ctx, input.Name)
		if err != nil {
			return nil, fmt.Errorf("repo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, except.ConflictError("category already exists")
		}

		categoryID := repo.GetNextID()

		category := entity.NewCategory(entity.CategoryParams{
			ID:              categoryID,
			Name:            input.Name,
			Icon:            input.Icon,
			CategoryGroupID: input.CategoryGroupID,
		})

		if err := repo.Store(ctx, category); err != nil {
			return nil, fmt.Errorf("repo.Store: %w", err)
		}

		return category, nil
	}
}
