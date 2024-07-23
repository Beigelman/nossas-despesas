package category

import (
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
	"time"
)

type GroupID struct{ Value int }

type Group struct {
	ddd.Entity[GroupID]
	Name string
	Icon string
}

type GroupAttributes struct {
	ID   GroupID
	Name string
	Icon string
}

func NewGroup(p GroupAttributes) *Group {
	return &Group{
		Entity: ddd.Entity[GroupID]{
			ID:        p.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Name: p.Name,
		Icon: p.Icon,
	}
}
