package userrepo

import (
	"database/sql"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToEntity(t *testing.T) {
	// Test with valid input
	model := UserModel{
		ID:             1,
		Name:           "Test User",
		Email:          "test@example.com",
		ProfilePicture: sql.NullString{String: "test.jpg", Valid: true},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DeletedAt:      sql.NullTime{Time: time.Now(), Valid: true},
		Version:        1,
	}
	user := toEntity(model)
	assert.NotNil(t, user)
	assert.Equal(t, entity.UserID{Value: 1}, user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "test.jpg", *user.ProfilePicture)
	assert.NotNil(t, user.DeletedAt)
	assert.Equal(t, 1, user.Version)

	// Test with invalid input
	model = UserModel{
		ID:             2,
		Name:           "Test User 2",
		Email:          "test2@example.com",
		ProfilePicture: sql.NullString{String: "", Valid: false},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DeletedAt:      sql.NullTime{Time: time.Time{}, Valid: false},
		Version:        2,
	}
	user = toEntity(model)
	assert.NotNil(t, user)
	assert.Equal(t, entity.UserID{Value: 2}, user.ID)
	assert.Equal(t, "Test User 2", user.Name)
	assert.Equal(t, "test2@example.com", user.Email)
	assert.Nil(t, user.ProfilePicture)
	assert.Nil(t, user.DeletedAt)
	assert.Equal(t, 2, user.Version)
}
