package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	mockrepository "github.com/Beigelman/ludaapi/internal/tests/mocks/repository"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateCategory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := mockrepository.NewMockCategoryRepository(t)
	existingCategory := entity.NewCategory(entity.CategoryParams{
		ID:              entity.CategoryID{Value: 1},
		Name:            "My test category",
		Icon:            "icon1",
		CategoryGroupID: entity.CategoryGroupID{Value: 1},
	})

	input := usecase.CreateCategoryInput{
		Name:            "My test category",
		Icon:            "icon1",
		CategoryGroupID: entity.CategoryGroupID{Value: 1},
	}

	useCase := usecase.NewCreateCategory(repo)

	t.Run("getByName returns error", func(t *testing.T) {
		repo.EXPECT().GetByName(ctx, existingCategory.Name).Return(nil, errors.New("test error")).Once()
		group, err := useCase(ctx, input)
		assert.Errorf(t, err, "repo.GetByName: test error")
		assert.Nil(t, group)
	})

	t.Run("category already exists", func(t *testing.T) {
		repo.EXPECT().GetByName(ctx, existingCategory.Name).Return(existingCategory, nil).Once()
		group, err := useCase(ctx, input)
		assert.Errorf(t, err, "category already exists")
		assert.Nil(t, group)
	})

	t.Run("Store returns error", func(t *testing.T) {
		repo.EXPECT().GetByName(ctx, existingCategory.Name).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		repo.EXPECT().GetNextID().Return(entity.CategoryID{Value: 1}).Once()
		group, err := useCase(ctx, input)
		assert.Errorf(t, err, "repo.Store: test error")
		assert.Nil(t, group)
	})

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetByName(ctx, mock.Anything).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		repo.EXPECT().GetNextID().Return(entity.CategoryID{Value: 1}).Once()
		category, err := useCase(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, entity.CategoryID{Value: 1}, category.ID)
		assert.Equal(t, "My test category", category.Name)
	})

}
