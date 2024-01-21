package config

import (
	"fmt"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/env"
)

type Db struct {
	Host          string        `mapstructure:"DB_HOST"`
	Port          string        `mapstructure:"DB_PORT"`
	Name          string        `mapstructure:"DB_NAME"`
	User          string        `mapstructure:"DB_USER"`
	Password      string        `mapstructure:"DB_PASSWORD"`
	Type          string        `mapstructure:"DB_TYPE"`
	MigrationPath string        `mapstructure:"DB_MIGRATION_PATH"`
	MaxIdleConns  int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns  int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxLifeTime   time.Duration `mapstructure:"DB_MAX_LIFE_TIME"`
	MaxIdleTime   time.Duration `mapstructure:"DB_MAX_IDLE_TIME"`
}

type Config struct {
	Env         env.Environment `mapstructure:"ENV"`
	ServiceName string          `mapstructure:"SERVICE_NAME"`
	Port        string          `mapstructure:"PORT"`
	LogLevel    string          `mapstructure:"LOG_LEVEL"`
	Db          Db              `mapstructure:",squash"`
}

func New(env env.Environment) Config {
	return Config{
		Env:  env,
		Port: "8080",
		Db: Db{
			Type:         "postgres",
			MaxOpenConns: 4,
		},
	}
}

func NewTestConfig(dbPort, dbHost, dbType string) Config {
	return Config{
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
			MigrationPath: "file:///Users/danielbeigelman/mydev/go-luda/server/database/migrations",
			MaxOpenConns:  4,
		},
	}
}

func (c Config) DBConnectionURL() string {
	if c.Db.Type == "postgres" {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", c.Db.Host, c.Db.Port, c.Db.User, c.Db.Name, c.Db.Password)
	} else if c.Db.Type == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Db.User, c.Db.Password, c.Db.Host, c.Db.Port, c.Db.Name)
	}

	panic("Invalid database type")
}
