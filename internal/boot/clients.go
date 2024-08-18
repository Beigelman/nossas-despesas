package boot

import (
	"context"
	"log/slog"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/email"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/jwt"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/pubsub"
	"github.com/Beigelman/nossas-despesas/internal/shared/service"
)

var ClientsModule = eon.NewModule("Clients", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, func(cfg *config.Config) service.TokenProvider {
		return jwt.NewJWTProvider(cfg.JWTSecret)
	})
	di.Provide(c, func(cfg *config.Config) service.EmailProvider {
		if cfg.Env == env.Development {
			return email.NewMailTrapEmailProvider(cfg.Mail.ApiKey)
		}
		return email.NewResendEmailProvider(cfg.Mail.ApiKey)
	})
	di.Provide(c, pubsub.NewSqlPublisher)
	di.Provide(c, pubsub.NewSqlSubscriber)

	lc.OnDisposing(eon.HookOrders.APPEND, func() error {
		if publisher := di.Resolve[pubsub.Publisher](c); publisher != nil {
			slog.InfoContext(ctx, "Closing publisher connection")
			return publisher.Close()
		}
		return nil
	})

	lc.OnDisposing(eon.HookOrders.APPEND, func() error {
		if subscriber := di.Resolve[pubsub.Subscriber](c); subscriber != nil {
			slog.InfoContext(ctx, "Closing subscriber connection")
			return subscriber.Close()
		}
		return nil
	})
})
