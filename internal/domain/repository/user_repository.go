package repository

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type UserRepository interface {
	ddd.Repository[entity.UserID, entity.User]
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
