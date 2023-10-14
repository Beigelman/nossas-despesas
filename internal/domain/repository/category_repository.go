package repository

import (
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type CategoryRepository interface {
	ddd.Repository[entity.CategoryID, entity.Category]
}
