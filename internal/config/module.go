package config

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/pkg/logger"
	"log/slog"
	"os"
)

const configPath = "./internal/config/config.yml"

var Module = eon.NewModule("Config", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, func() (*Config, error) {
		environment, err := env.Parse(os.Getenv("ENV"))
		if err != nil {
			return nil, fmt.Errorf("env.Parse: %w", err)
		}

		cfg := New(environment)
		cfg.SetConfigPath(configPath)
		if err := cfg.LoadConfig(); err != nil {
			return nil, fmt.Errorf("cfg.LoadConfig: %w", err)
		}

		cfg.ServiceName = info.ServiceName

		return &cfg, nil
	})

	lc.OnBooted(eon.HookOrders.PREPEND, func() error {
		cfg := di.Resolve[*Config](c)
		if cfg.Env == env.Development {
			slog.SetDefault(logger.NewDevelopment(cfg.LogLevel))
		} else {
			slog.SetDefault(logger.NewProduction(cfg.LogLevel))
		}

		return nil
	})
})
