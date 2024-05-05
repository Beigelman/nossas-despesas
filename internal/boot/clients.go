package boot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/domain/service"
	"github.com/Beigelman/nossas-despesas/internal/infra/email"
	"github.com/Beigelman/nossas-despesas/internal/infra/jwt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var ClientsModule = eon.NewModule("Clients", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	var dbClient db.Database
	di.Provide(c, func(cfg *config.Config) service.TokenProvider {
		return jwt.NewJWTProvider(cfg.JWTSecret)
	})
	di.Provide(c, func(cfg *config.Config) db.Database {
		dbClient = db.New(cfg)
		return dbClient
	})

	di.Provide(c, func(cfg *config.Config) service.EmailProvider {
		if cfg.Env == env.Development {
			return email.NewMailTrapEmailProvider(cfg.Mail.ApiKey)
		}

		return email.NewResendEmailProvider(cfg.Mail.ApiKey)
	})

	lc.OnDisposing(eon.HookOrders.PREPEND, func() error {
		if dbClient != nil {
			slog.Info("Closing db connection")
			if err := dbClient.Close(); err != nil {
				return fmt.Errorf("dbClient.Close: %w", err)
			}
		}

		return nil
	})
})
