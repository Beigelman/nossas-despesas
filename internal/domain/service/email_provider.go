package service

import (
	"context"
	vo "github.com/Beigelman/nossas-despesas/internal/domain/valueobject"
)

type EmailProvider interface {
	Send(ctx context.Context, email vo.Email) error
}
