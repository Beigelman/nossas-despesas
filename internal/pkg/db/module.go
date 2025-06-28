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
	di.Provide(c, func(cfg *config.Config) (*Client, error) {
		dbClient, err := NewClient(
			cfg.DBConnectionString(),
			WithConnMaxIdleTime(cfg.Db.MaxIdleTime),
			WithConnMaxLifeTime(cfg.Db.MaxLifeTime),
			WithMaxIdleConns(cfg.Db.MaxIdleConns),
			WithMaxOpenConns(cfg.Db.MaxOpenConns),
			WithMigrationPath(cfg.Db.MigrationPath),
		)
		if err != nil {
			return nil, fmt.Errorf("db.New: %w", err)
		}

		return dbClient, nil
	})

	lc.OnDisposing(eon.HookOrders.PREPEND, func() error {
		slog.Info("Closing db connection")
		dbClient := di.Resolve[*Client](c)
		if err := dbClient.Close(); err != nil {
			return fmt.Errorf("dbClient.Close: %w", err)
		}

		return nil
	})
})
