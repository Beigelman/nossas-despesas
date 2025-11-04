package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/category/usecase"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCategory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := mocks.NewMockcategoryRepository(t)
	existingCategory := category.New(category.Attributes{
		ID:              category.ID{Value: 1},
		Name:            "My test category",
		Icon:            "icon1",
		CategoryGroupID: category.GroupID{Value: 1},
	})

	input := usecase.CreateCategoryInput{
		Name:            "My test category",
		Icon:            "icon1",
		CategoryGroupID: category.GroupID{Value: 1},
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
		repo.EXPECT().GetNextID().Return(category.ID{Value: 1}).Once()
		group, err := useCase(ctx, input)
		assert.Errorf(t, err, "repo.Store: test error")
		assert.Nil(t, group)
	})

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetByName(ctx, mock.Anything).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		repo.EXPECT().GetNextID().Return(category.ID{Value: 1}).Once()
		cat, err := useCase(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, category.ID{Value: 1}, cat.ID)
		assert.Equal(t, "My test category", cat.Name)
	})
}
