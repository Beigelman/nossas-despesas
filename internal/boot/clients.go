package boot

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/email"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/jwt"
	"github.com/Beigelman/nossas-despesas/internal/shared/service"
)

var ClientsModule = eon.NewModule("Clients", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, func(cfg *config.Config) *jwt.Provider {
		return jwt.NewJWTProvider(cfg.JWTSecret)
	})
	di.Provide(c, func(cfg *config.Config) service.EmailProvider {
		if cfg.Env == env.Development {
			return email.NewMailTrapEmailProvider(cfg.Mail.ApiKey)
		}

		return email.NewResendEmailProvider(cfg.Mail.ApiKey)
	})
})
