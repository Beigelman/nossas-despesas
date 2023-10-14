package repository

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/domain/entity"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type UserRepository interface {
	ddd.Repository[entity.UserID, entity.User]
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
