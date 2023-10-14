package categoryrepo

import (
	"database/sql"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"testing"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
	"github.com/stretchr/testify/assert"
)

func TestToEntity(t *testing.T) {
	// Test with valid input
	model := CategoryModel{
		ID:        1,
		Name:      "Test Category",
		Icon:      "test.png",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Version:   1,
	}
	category := toEntity(model)
	assert.NotNil(t, category)
	assert.Equal(t, entity.CategoryID{Value: 1}, category.ID)
	assert.Equal(t, "Test Category", category.Name)
	assert.Equal(t, "test.png", category.Icon)
	assert.NotNil(t, category.DeletedAt)
	assert.Equal(t, 1, category.Version)

	// Test with invalid input
	model = CategoryModel{
		ID:        2,
		Name:      "Test Category 2",
		Icon:      "test2.png",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: sql.NullTime{Time: time.Time{}, Valid: false},
		Version:   2,
	}
	category = toEntity(model)
	assert.NotNil(t, category)
	assert.Equal(t, entity.CategoryID{Value: 2}, category.ID)
	assert.Equal(t, "Test Category 2", category.Name)
	assert.Equal(t, "test2.png", category.Icon)
	assert.Nil(t, category.DeletedAt)
	assert.Equal(t, 2, category.Version)
}

func TestToModel(t *testing.T) {
	// Test with valid input
	deletedAt := time.Now()
	category := &entity.Category{
		Entity: ddd.Entity[entity.CategoryID]{
			ID:        entity.CategoryID{Value: 1},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: &deletedAt,
			Version:   1,
		},
		Name: "Test Category",
		Icon: "test.png",
	}
	model := toModel(category)
	assert.NotNil(t, model)
	assert.Equal(t, 1, model.ID)
	assert.Equal(t, "Test Category", model.Name)
	assert.Equal(t, "test.png", model.Icon)
	assert.True(t, model.CreatedAt.Before(time.Now()))
	assert.True(t, model.UpdatedAt.Before(time.Now()))
	assert.True(t, model.DeletedAt.Valid)
	assert.Equal(t, 1, model.Version)

	// Test with invalid input
	category = &entity.Category{
		Entity: ddd.Entity[entity.CategoryID]{
			ID:        entity.CategoryID{Value: 2},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
			Version:   2,
		},
		Name: "Test Category 2",
		Icon: "test2.png",
	}
	model = toModel(category)
	assert.NotNil(t, model)
	assert.Equal(t, 2, model.ID)
	assert.Equal(t, "Test Category 2", model.Name)
	assert.Equal(t, "test2.png", model.Icon)
	assert.True(t, model.CreatedAt.Before(time.Now()))
	assert.True(t, model.UpdatedAt.Before(time.Now()))
	assert.False(t, model.DeletedAt.Valid)
	assert.Equal(t, 2, model.Version)
}
