package user

import (
	"context"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type Repository interface {
	ddd.Repository[ID, User]
	GetByEmail(ctx context.Context, email string) (*User, error)
}
