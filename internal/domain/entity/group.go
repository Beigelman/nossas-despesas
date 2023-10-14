package entity

import (
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type GroupID struct{ Value int }

type Group struct {
	ddd.Entity[GroupID]
	Name string
}

type GroupParams struct {
	ID   GroupID
	Name string
}

func NewGroup(p GroupParams) *Group {
	return &Group{
		Entity: ddd.Entity[GroupID]{
			ID:        p.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Name: p.Name,
	}
}

func (g *Group) SetName(name string) {
	g.Name = name
}
