package ddd

import "context"

type Repository[ID any, Entity any] interface {
	GetNextID() ID
	GetByID(ctx context.Context, id ID) (*Entity, error)
	Store(ctx context.Context, entity *Entity) error
}
