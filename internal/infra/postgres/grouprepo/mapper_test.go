package grouprepo

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
	model := GroupModel{
		ID:        1,
		Name:      "Test Group",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Version:   1,
	}
	group := toEntity(model)
	assert.NotNil(t, group)
	assert.Equal(t, entity.GroupID{Value: 1}, group.ID)
	assert.Equal(t, "Test Group", group.Name)
	assert.NotNil(t, group.DeletedAt)
	assert.Equal(t, 1, group.Version)

	// Test with invalid input
	model = GroupModel{
		ID:        2,
		Name:      "Test Group 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: sql.NullTime{Time: time.Time{}, Valid: false},
		Version:   2,
	}
	group = toEntity(model)
	assert.NotNil(t, group)
	assert.Equal(t, entity.GroupID{Value: 2}, group.ID)
	assert.Equal(t, "Test Group 2", group.Name)
	assert.Nil(t, group.DeletedAt)
	assert.Equal(t, 2, group.Version)
}

func TestToModel(t *testing.T) {
	// Test with valid input
	deletedAt := time.Now()
	group := &entity.Group{
		Entity: ddd.Entity[entity.GroupID]{
			ID:        entity.GroupID{Value: 1},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: &deletedAt,
			Version:   1,
		},
		Name: "Test Group",
	}
	model := toModel(group)
	assert.NotNil(t, model)
	assert.Equal(t, 1, model.ID)
	assert.Equal(t, "Test Group", model.Name)
	assert.True(t, model.CreatedAt.Before(time.Now()))
	assert.True(t, model.UpdatedAt.Before(time.Now()))
	assert.True(t, model.DeletedAt.Valid)
	assert.Equal(t, 1, model.Version)

	// Test with invalid input
	group = &entity.Group{
		Entity: ddd.Entity[entity.GroupID]{
			ID:        entity.GroupID{Value: 2},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
			Version:   2,
		},
		Name: "Test Group 2",
	}
	model = toModel(group)
	assert.NotNil(t, model)
	assert.Equal(t, 2, model.ID)
	assert.Equal(t, "Test Group 2", model.Name)
	assert.True(t, model.CreatedAt.Before(time.Now()))
	assert.True(t, model.UpdatedAt.Before(time.Now()))
	assert.False(t, model.DeletedAt.Valid)
	assert.Equal(t, 2, model.Version)
}
