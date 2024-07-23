package auth

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type Repository interface {
	ddd.Repository[ID, Auth]
	GetByEmail(ctx context.Context, email string, authType Type) (*Auth, error)
}
