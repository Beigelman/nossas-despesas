package boot

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/email"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/jwt"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/pubsub"
	"github.com/Beigelman/nossas-despesas/internal/shared/service"
	"github.com/ThreeDotsLabs/watermill/message"
	"log/slog"
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
	di.Provide(c, func(db db.Database) (message.Publisher, error) {
		return pubsub.NewSqlPublisher(db.Client())
	})
	di.Provide(c, func(db db.Database) (message.Subscriber, error) {
		return pubsub.NewSqlSubscriber(db.Client())
	})

	lc.OnDisposing(eon.HookOrders.APPEND, func() error {
		if publisher := di.Resolve[message.Publisher](c); publisher != nil {
			slog.InfoContext(ctx, "Closing publisher connection")
			return publisher.Close()
		}
		return nil
	})

	lc.OnDisposing(eon.HookOrders.APPEND, func() error {
		if subscriber := di.Resolve[message.Subscriber](c); subscriber != nil {
			slog.InfoContext(ctx, "Closing subscriber connection")
			return subscriber.Close()
		}
		return nil
	})
})
