package repository

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type AuthRepository interface {
	ddd.Repository[entity.AuthID, entity.Auth]
	GetByEmail(ctx context.Context, email string, authType entity.AuthType) (*entity.Auth, error)
}
