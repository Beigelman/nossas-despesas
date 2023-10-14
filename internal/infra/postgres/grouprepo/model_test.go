package grouprepo

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGroupModel(t *testing.T) {
	// Test fields
	group := GroupModel{
		ID:        1,
		Name:      "Test Group",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Version:   1,
	}
	assert.Equal(t, 1, group.ID)
	assert.Equal(t, "Test Group", group.Name)
	assert.True(t, group.CreatedAt.Before(time.Now()))
	assert.True(t, group.UpdatedAt.Before(time.Now()))
	assert.True(t, group.DeletedAt.Valid)
	assert.Equal(t, 1, group.Version)
}
