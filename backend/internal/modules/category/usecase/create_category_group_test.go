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

func TestCreateCategoryGroup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := mocks.NewMockcategoryGroupRepository(t)
	existingCategory := category.NewGroup(category.GroupAttributes{
		ID:   category.GroupID{Value: 1},
		Name: "My test category",
		Icon: "icon1",
	})

	input := usecase.CreateCategoryGroupInput{
		Name: "My test category",
		Icon: "icon1",
	}

	useCase := usecase.NewCreateCategoryGroup(repo)

	t.Run("getByName returns error", func(t *testing.T) {
		repo.EXPECT().GetByName(ctx, existingCategory.Name).Return(nil, errors.New("test error")).Once()
		group, err := useCase(ctx, input)
		assert.Errorf(t, err, "repo.GetByName: test error")
		assert.Nil(t, group)
	})

	t.Run("group category already exists", func(t *testing.T) {
		repo.EXPECT().GetByName(ctx, existingCategory.Name).Return(existingCategory, nil).Once()
		group, err := useCase(ctx, input)
		assert.Errorf(t, err, "category already exists")
		assert.Nil(t, group)
	})

	t.Run("Store returns error", func(t *testing.T) {
		repo.EXPECT().GetByName(ctx, existingCategory.Name).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		repo.EXPECT().GetNextID().Return(category.GroupID{Value: 1}).Once()
		group, err := useCase(ctx, input)
		assert.Errorf(t, err, "repo.Store: test error")
		assert.Nil(t, group)
	})

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetByName(ctx, mock.Anything).Return(nil, nil).Once()
		repo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		repo.EXPECT().GetNextID().Return(category.GroupID{Value: 1}).Once()
		categoryGroup, err := useCase(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, categoryGroup)
		assert.Equal(t, category.GroupID{Value: 1}, categoryGroup.ID)
		assert.Equal(t, "My test category", categoryGroup.Name)
	})
}
