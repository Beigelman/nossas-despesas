package config

import (
	"fmt"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/env"
)

type Db struct {
	Host         string        `mapstructure:"DB_HOST"`
	Port         string        `mapstructure:"DB_PORT"`
	Name         string        `mapstructure:"DB_NAME"`
	User         string        `mapstructure:"DB_USER"`
	Password     string        `mapstructure:"DB_PASSWORD"`
	Type         string        `mapstructure:"DB_TYPE"`
	MaxIdleConns int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxLifeTime  time.Duration `mapstructure:"DB_MAX_LIFE_TIME"`
	MaxIdleTime  time.Duration `mapstructure:"DB_MAX_IDLE_TIME"`
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

func NewTestConfig(dbPort string, dbHost string) Config {
	return Config{
		ServiceName: "test-luda-api",
		Env:         "test",
		Port:        "8080",
		Db: Db{
			Type:         "postgres",
			Host:         dbHost,
			User:         "root",
			Password:     "root",
			Name:         "test",
			Port:         dbPort,
			MaxOpenConns: 4,
		},
	}
}

func (c Config) DBConnectionURL() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", c.Db.Host, c.Db.User, c.Db.Password, c.Db.Name, c.Db.Port)
}
