package boot

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/domain/service"
	"github.com/Beigelman/ludaapi/internal/infra/jwt"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/Beigelman/ludaapi/internal/pkg/eon"
	"google.golang.org/api/option"
	"log/slog"
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
	di.Provide(c, func() (*auth.Client, error) {
		opt := option.WithCredentialsFile("./firebaseServiceAccount.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			return nil, fmt.Errorf("firebase.NewApp: %w", err)
		}
		return app.Auth(ctx)
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
