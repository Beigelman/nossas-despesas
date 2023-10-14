package categoryrepo

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCategoryModel(t *testing.T) {
	// Test fields
	category := CategoryModel{
		ID:        1,
		Name:      "Test Category",
		Icon:      "test.png",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Version:   1,
	}
	assert.Equal(t, 1, category.ID)
	assert.Equal(t, "Test Category", category.Name)
	assert.Equal(t, "test.png", category.Icon)
	assert.True(t, category.CreatedAt.Before(time.Now()))
	assert.True(t, category.UpdatedAt.Before(time.Now()))
	assert.True(t, category.DeletedAt.Valid)
	assert.Equal(t, 1, category.Version)

}
