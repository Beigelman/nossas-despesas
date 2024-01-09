package repository

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/domain/entity"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type AuthRepository interface {
	ddd.Repository[entity.AuthID, entity.Auth]
	GetByEmail(ctx context.Context, email string, authType entity.AuthType) (*entity.Auth, error)
}
