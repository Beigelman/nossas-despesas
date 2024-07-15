package user

import (
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ID struct{ Value int }

type User struct {
	ddd.Entity[ID]
	Name           string
	Email          string
	ProfilePicture *string
	GroupID        *entity.GroupID
}

type Attributes struct {
	ID             ID
	Name           string
	Email          string
	ProfilePicture *string
	GroupID        *entity.GroupID
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

func (u *User) AssignGroup(g entity.GroupID) {
	u.GroupID = &g
}
