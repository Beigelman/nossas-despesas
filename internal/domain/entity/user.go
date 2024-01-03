package entity

import (
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type UserID struct{ Value int }

type User struct {
	ddd.Entity[UserID]
	Name             string
	Email            string
	ProfilePicture   *string
	GroupID          *GroupID
	AuthenticationID *string
}

type UserParams struct {
	ID               UserID
	Name             string
	Email            string
	ProfilePicture   *string
	GroupID          *GroupID
	AuthenticationID *string
}

func NewUser(p UserParams) *User {
	return &User{
		Entity: ddd.Entity[UserID]{
			ID:        p.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Name:             p.Name,
		Email:            p.Email,
		ProfilePicture:   p.ProfilePicture,
		GroupID:          p.GroupID,
		AuthenticationID: p.AuthenticationID,
	}
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) AssignGroup(g GroupID) {
	u.GroupID = &g
}
