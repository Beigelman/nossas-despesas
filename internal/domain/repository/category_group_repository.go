package repository

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type CategoryGroupRepository interface {
	ddd.Repository[entity.CategoryGroupID, entity.CategoryGroup]
	GetByName(ctx context.Context, name string) (*entity.CategoryGroup, error)
}
