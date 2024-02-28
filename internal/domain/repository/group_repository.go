package repository

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type GroupRepository interface {
	ddd.Repository[entity.GroupID, entity.Group]
	GetByName(ctx context.Context, name string) (*entity.Group, error)
}
