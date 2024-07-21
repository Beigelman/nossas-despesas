package category

import (
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ID struct{ Value int }

type Category struct {
	ddd.Entity[ID]
	Name            string
	Icon            string
	GroupCategoryID GroupID
}

type Attributes struct {
	ID              ID
	Name            string
	Icon            string
	CategoryGroupID GroupID
}

func NewCategory(p Attributes) *Category {
	return &Category{
		Entity: ddd.Entity[ID]{
			ID:        p.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		Name:            p.Name,
		Icon:            p.Icon,
		GroupCategoryID: p.CategoryGroupID,
	}
}
