package category

import (
	"context"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
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

type GroupRepository interface {
	ddd.Repository[GroupID, Group]
	GetByName(ctx context.Context, name string) (*Group, error)
}
