package user

import (
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ID struct{ Value int }

type User struct {
	ddd.Entity[ID]
	Name           string
	Email          string
	ProfilePicture *string
	GroupID        *group.ID
	Flags          []Flag
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
		Flags:          []Flag{},
	}
}

func (u *User) HasFlag(flag Flag) bool {
	for _, f := range u.Flags {
		if f == flag {
			return true
		}
	}

	return false
}

func (u *User) AddFlag(flag Flag) {
	u.Flags = append(u.Flags, flag)
}

func (u *User) RemoveFlag(flag Flag) {
	filtredFlags := u.Flags[:0]
	for _, f := range u.Flags {
		if f == flag {
			continue
		}

		filtredFlags = append(filtredFlags, f)
	}

	u.Flags = filtredFlags
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) AssignGroup(g group.ID) {
	u.GroupID = &g
}
