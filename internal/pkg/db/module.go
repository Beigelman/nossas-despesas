package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Beigelman/nossas-despesas/internal/pkg/config"

	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var Module = eon.NewModule("Database", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, func(cfg *config.Config) (Database, error) {
		dbClient, err := New(cfg)
		if err != nil {
			return nil, fmt.Errorf("db.New: %w", err)
		}

		return dbClient, nil
	})

	lc.OnDisposing(eon.HookOrders.PREPEND, func() error {
		slog.Info("Closing db connection")
		dbClient := di.Resolve[Database](c)
		if err := dbClient.Close(); err != nil {
			return fmt.Errorf("dbClient.Close: %w", err)
		}

		return nil
	})
})
