package group

import (
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ID struct{ Value int }

type Group struct {
	ddd.Entity[ID]
	Name string
}

type Attributes struct {
	ID   ID
	Name string
}

func New(p Attributes) *Group {
	return &Group{
		Entity: ddd.Entity[ID]{
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
