package db

import (
	"context"
	"fmt"
	"time"

	nossasdespesas "github.com/Beigelman/nossas-despesas"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var Module = eon.NewModule("Database", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, func(cfg *nossasdespesas.Config) (*Client, error) {
		dbClient, err := NewClient(
			cfg.DBConnectionString(),
			WithConnMaxIdleTime(time.Duration(cfg.Db.MaxIdleTimeMinutes)*time.Minute),
			WithConnMaxLifeTime(time.Duration(cfg.Db.MaxLifeTimeMinutes)*time.Minute),
			WithMaxIdleConns(cfg.Db.MaxIdleConns),
			WithMaxOpenConns(cfg.Db.MaxOpenConns),
		)
		if err != nil {
			return nil, fmt.Errorf("db.New: %w", err)
		}

		return dbClient, nil
	})

	lc.OnDisposing(eon.HookOrders.PREPEND, func() error {
		dbClient := di.Resolve[*Client](c)
		if err := dbClient.Close(); err != nil {
			return fmt.Errorf("dbClient.Close: %w", err)
		}

		return nil
	})
})
