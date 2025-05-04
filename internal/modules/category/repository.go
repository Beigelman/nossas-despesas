package category

import (
	"context"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type Repository interface {
	ddd.Repository[ID, Category]
	GetByName(ctx context.Context, name string) (*Category, error)
}

type GroupRepository interface {
	ddd.Repository[GroupID, Group]
	GetByName(ctx context.Context, name string) (*Group, error)
}
