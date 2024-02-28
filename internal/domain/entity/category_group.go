package entity

import (
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
	"time"
)

type CategoryGroupID struct{ Value int }

type CategoryGroup struct {
	ddd.Entity[CategoryGroupID]
	Name string
	Icon string
}

type CategoryGroupParams struct {
	ID   CategoryGroupID
	Name string
	Icon string
}

func NewCategoryGroup(p CategoryGroupParams) *CategoryGroup {
	return &CategoryGroup{
		Entity: ddd.Entity[CategoryGroupID]{
			ID:        p.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Name: p.Name,
		Icon: p.Icon,
	}
}
