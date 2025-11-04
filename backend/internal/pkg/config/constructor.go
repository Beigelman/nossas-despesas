package config

import (
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/spf13/viper"
)

type Db struct {
	Host             string        `mapstructure:"DB_HOST"`
	Port             string        `mapstructure:"DB_PORT"`
	Name             string        `mapstructure:"DB_NAME"`
	User             string        `mapstructure:"DB_USER"`
	Password         string        `mapstructure:"DB_PASSWORD"`
	ConnectionString string        `mapstructure:"DB_CONNECTION_STRING"`
	MigrationPath    string        `mapstructure:"DB_MIGRATION_PATH"`
	MaxIdleConns     int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns     int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxLifeTime      time.Duration `mapstructure:"DB_MAX_LIFE_TIME"`
	MaxIdleTime      time.Duration `mapstructure:"DB_MAX_IDLE_TIME"`
}

type Mail struct {
	SandBoxID string `mapstructure:"MAIL_SANDBOX_ID"`
	ApiKey    string `mapstructure:"MAIL_API_KEY"`
}

type Config struct {
	loader      *viper.Viper
	Env         env.Environment `mapstructure:"ENV"`
	ServiceName string          `mapstructure:"SERVICE_NAME"`
	Port        string          `mapstructure:"PORT"`
	LogLevel    string          `mapstructure:"LOG_LEVEL"`
	JWTSecret   string          `mapstructure:"JWT_SECRET"`
	Mail        Mail            `mapstructure:",squash"`
	Db          Db              `mapstructure:",squash"`
	SentryDsn   string          `mapstructure:"SENTRY_DSN"`
}

func New(env env.Environment) Config {
	return Config{
		loader:   viper.New(),
		Env:      env,
		Port:     "8080",
		LogLevel: "INFO",
		Db: Db{
			MaxOpenConns: 4,
		},
	}
}
