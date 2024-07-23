package boot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var DatabaseModule = eon.NewModule("Database", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, func(cfg *config.Config) (db.Database, error) {
		dbClient, err := db.New(cfg)
		if err != nil {
			return nil, fmt.Errorf("db.New: %w", err)
		}

		return dbClient, nil
	})

	lc.OnDisposing(eon.HookOrders.PREPEND, func() error {
		slog.Info("Closing db connection")
		dbClient := di.Resolve[db.Database](c)
		if err := dbClient.Close(); err != nil {
			return fmt.Errorf("dbClient.Close: %w", err)
		}

		return nil
	})
})
