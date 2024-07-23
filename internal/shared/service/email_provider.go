package service

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/email"
)

type EmailProvider interface {
	Send(ctx context.Context, email email.Email) error
}
