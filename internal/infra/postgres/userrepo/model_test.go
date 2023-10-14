package userrepo_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Beigelman/ludaapi/internal/infra/postgres/userrepo"
	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	// Test fields
	user := userrepo.UserModel{
		ID:             1,
		Name:           "Test User",
		Email:          "test@example.com",
		ProfilePicture: sql.NullString{String: "test.jpg", Valid: true},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DeletedAt:      sql.NullTime{Time: time.Now(), Valid: true},
		Version:        1,
	}

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "test.jpg", user.ProfilePicture.String)
	assert.True(t, user.ProfilePicture.Valid)
	assert.True(t, user.CreatedAt.Before(time.Now()))
	assert.True(t, user.UpdatedAt.Before(time.Now()))
	assert.True(t, user.DeletedAt.Valid)
	assert.Equal(t, 1, user.Version)
}
