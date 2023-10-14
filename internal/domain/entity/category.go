package entity

import (
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type CategoryID struct{ Value int }

type Category struct {
	ddd.Entity[CategoryID]
	Name string
	Icon string
}

type CategoryParams struct {
	ID   CategoryID
	Name string
	Icon string
}

func NewCategory(p CategoryParams) *Category {
	return &Category{
		Entity: ddd.Entity[CategoryID]{
			ID:        p.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		Name: p.Name,
		Icon: p.Icon,
	}
}
