package boot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/Beigelman/ludaapi/internal/pkg/eon"
)

var DatabaseModule = eon.NewModule("Database", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	var dbClient db.Database

	di.Provide(c, func(cfg *config.Config) db.Database {
		dbClient = db.New(cfg)
		return dbClient
	})

	lc.OnDisposing(eon.HookOrders.PREPEND, func() error {
		slog.Info("Closing db connection")
		if err := dbClient.Close(); err != nil {
			return fmt.Errorf("dbClient.Close: %w", err)
		}

		return nil
	})
})
