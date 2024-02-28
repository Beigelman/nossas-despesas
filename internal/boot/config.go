package boot

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"log/slog"
	"os"
)

const configPath = "./internal/config/config.yml"

func LogLevelMap(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

var ConfigModule = eon.NewModule("Config", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, func() (*config.Config, error) {
		environment, err := env.Parse(os.Getenv("ENV"))
		if err != nil {
			return nil, fmt.Errorf("env.Parse: %w", err)
		}

		cfg := config.New(environment)
		cfg.SetConfigPath(configPath)
		if err := cfg.LoadConfig(); err != nil {
			return nil, fmt.Errorf("cfg.LoadConfig: %w", err)
		}

		cfg.ServiceName = info.ServiceName

		return &cfg, nil
	})

	lc.OnBooted(eon.HookOrders.PREPEND, func() error {
		cfg := di.Resolve[*config.Config](c)
		slog.SetDefault(slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: LogLevelMap(cfg.LogLevel),
				ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
					if a.Key == "level" {
						return slog.Attr{Key: "severity", Value: a.Value}
					}

					if a.Key == "msg" {
						return slog.Attr{Key: "message", Value: a.Value}
					}

					if a.Key == "time" {
						return slog.Attr{Key: "timestamp", Value: a.Value}
					}

					return a
				}})))
		return nil
	})
})
