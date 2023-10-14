package usecase

import (
	"context"
	"errors"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateGroup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	groupRepo := mocks.NewMockGroupRepository(t)
	existingGroup := entity.NewGroup(entity.GroupParams{
		ID:   entity.GroupID{Value: 1},
		Name: "My test group",
	})

	useCase := NewCreateGroup(groupRepo)

	t.Run("getByName returns error", func(t *testing.T) {
		groupRepo.EXPECT().GetByName(ctx, existingGroup.Name).Return(nil, errors.New("test error")).Once()
		group, err := useCase(ctx, existingGroup.Name)
		assert.Errorf(t, err, "repo.GetByName: test error")
		assert.Nil(t, group)
	})

	t.Run("group already exists", func(t *testing.T) {
		groupRepo.EXPECT().GetByName(ctx, existingGroup.Name).Return(existingGroup, nil).Once()
		group, err := useCase(ctx, existingGroup.Name)
		assert.Errorf(t, err, "group already exists")
		assert.Nil(t, group)
	})

	t.Run("Store returns error", func(t *testing.T) {
		groupRepo.EXPECT().GetByName(ctx, existingGroup.Name).Return(nil, nil).Once()
		groupRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
		groupRepo.EXPECT().GetNextID().Return(entity.GroupID{Value: 1}).Once()
		group, err := useCase(ctx, existingGroup.Name)
		assert.Errorf(t, err, "repo.Store: test error")
		assert.Nil(t, group)
	})

	t.Run("success", func(t *testing.T) {
		groupRepo.EXPECT().GetByName(ctx, mock.Anything).Return(nil, nil).Once()
		groupRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		groupRepo.EXPECT().GetNextID().Return(entity.GroupID{Value: 1}).Once()
		group, err := useCase(ctx, "my new group")
		assert.NoError(t, err)
		assert.NotNil(t, group)
		assert.Equal(t, entity.GroupID{Value: 1}, group.ID)
		assert.Equal(t, "my new group", group.Name)
	})

}
