package boot

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/Beigelman/ludaapi/internal/pkg/eon"
	"github.com/spf13/viper"
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
		v := viper.New()

		extension := filepath.Ext(configPath)                 // eg: .yml
		filename := filepath.Base(configPath)                 // eg: config.yml
		configName := strings.TrimSuffix(filename, extension) // eg: config
		v.SetConfigName(configName)                           // viper takes filename without extension

		if len(extension) > 1 {
			configType := extension[1:]
			v.SetConfigType(configType)
		}

		configDir := filepath.Dir(configPath) // eg: /app or .
		v.AddConfigPath(configDir)

		environment, err := env.Parse(os.Getenv("ENV"))
		if err != nil {
			return nil, fmt.Errorf("env.Parse: %w", err)
		}

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("viper.ReadInConfig: %w", err)
		}

		envViper := v.Sub(environment.String())
		envViper.SetDefault("PORT", "8080")
		envViper.SetDefault("LOG_LEVEL", "INFO")

		config := config.New(environment)
		err = envViper.Unmarshal(&config)
		if err != nil {
			return nil, fmt.Errorf("viper.Unmarshal: %w", err)
		}

		config.ServiceName = info.ServiceName

		return &config, nil
	})

	lc.OnBooted(eon.HookOrders.PREPEND, func() error {
		cfg := di.Resolve[*config.Config](c)
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: LogLevelMap(cfg.LogLevel)})))
		return nil
	})
})
