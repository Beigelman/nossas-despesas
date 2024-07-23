package user

import (
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ID struct{ Value int }

type User struct {
	ddd.Entity[ID]
	Name           string
	Email          string
	ProfilePicture *string
	GroupID        *group.ID
}

type Attributes struct {
	ID             ID
	Name           string
	Email          string
	ProfilePicture *string
	GroupID        *group.ID
}

func New(p Attributes) *User {
	return &User{
		Entity: ddd.Entity[ID]{
			ID:        p.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Name:           p.Name,
		Email:          p.Email,
		ProfilePicture: p.ProfilePicture,
		GroupID:        p.GroupID,
	}
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) AssignGroup(g group.ID) {
	u.GroupID = &g
}
