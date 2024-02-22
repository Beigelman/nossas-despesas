package config

import (
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/env"
)

type Db struct {
	Host             string        `mapstructure:"DB_HOST"`
	Port             string        `mapstructure:"DB_PORT"`
	Name             string        `mapstructure:"DB_NAME"`
	User             string        `mapstructure:"DB_USER"`
	Password         string        `mapstructure:"DB_PASSWORD"`
	Type             string        `mapstructure:"DB_TYPE"`
	ConnectionString string        `mapstructure:"DB_CONNECTION_STRING"`
	MigrationPath    string        `mapstructure:"DB_MIGRATION_PATH"`
	MaxIdleConns     int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns     int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxLifeTime      time.Duration `mapstructure:"DB_MAX_LIFE_TIME"`
	MaxIdleTime      time.Duration `mapstructure:"DB_MAX_IDLE_TIME"`
}

type Config struct {
	loader      *viper.Viper
	Env         env.Environment `mapstructure:"ENV"`
	ServiceName string          `mapstructure:"SERVICE_NAME"`
	Port        string          `mapstructure:"PORT"`
	LogLevel    string          `mapstructure:"LOG_LEVEL"`
	JWTSecret   string          `mapstructure:"JWT_SECRET"`
	Db          Db              `mapstructure:",squash"`
}

func New(env env.Environment) Config {
	return Config{
		loader: viper.New(),
		Env:    env,
		Port:   "8080",
		Db: Db{
			Type:         "postgres",
			MaxOpenConns: 4,
		},
	}
}

func NewTestConfig(dbPort, dbHost, dbType string) Config {
	return Config{
		loader:      viper.New(),
		ServiceName: "test-luda-api",
		Env:         "test",
		Port:        "8080",
		Db: Db{
			Host:          dbHost,
			Port:          dbPort,
			Name:          "test",
			User:          "root",
			Password:      "root",
			Type:          dbType,
			MigrationPath: os.Getenv("MIGRATION_PATH"),
			MaxOpenConns:  4,
		},
	}
}
