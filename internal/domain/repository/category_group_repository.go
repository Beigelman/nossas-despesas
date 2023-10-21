package repository

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type CategoryGroupRepository interface {
	ddd.Repository[entity.CategoryGroupID, entity.CategoryGroup]
	GetByName(ctx context.Context, name string) (*entity.CategoryGroup, error)
}
