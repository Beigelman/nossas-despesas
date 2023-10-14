package repository

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/domain/entity"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type GroupRepository interface {
	ddd.Repository[entity.GroupID, entity.Group]
	GetByName(ctx context.Context, name string) (*entity.Group, error)
}
