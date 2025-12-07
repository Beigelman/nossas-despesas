package nossasdespesas

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"

	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
)

type Db struct {
	Host               string `env:"DB_HOST"`
	Port               string `env:"DB_PORT"`
	Name               string `env:"DB_NAME"`
	User               string `env:"DB_USER"`
	Password           string `env:"DB_PASSWORD"`
	ConnectionString   string `env:"DB_CONNECTION_STRING"`
	MigrationPath      string `env:"DB_MIGRATION_PATH"`
	MaxIdleConns       int    `env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns       int    `env:"DB_MAX_OPEN_CONNS"`
	MaxLifeTimeMinutes int    `env:"DB_MAX_LIFE_TIME_MINUTES"`
	MaxIdleTimeMinutes int    `env:"DB_MAX_IDLE_TIME_MINUTES"`
}

type Mail struct {
	SandBoxID string `env:"MAIL_SANDBOX_ID"`
	ApiKey    string `env:"MAIL_API_KEY"`
}

type Config struct {
	Env         env.Environment
	ServiceName string `env:"SERVICE_NAME"`
	Port        string `env:"PORT"`
	LogLevel    string `env:"LOG_LEVEL"`
	JWTSecret   string `env:"JWT_SECRET"`
	SentryDsn   string `env:"SENTRY_DSN"`
	PredictURL  string `env:"PREDICT_URL"`
	Mail        Mail
	Db          Db
}

func NewConfig(environment env.Environment) (Config, error) {
	cfg := Config{Env: environment}

	// Create config builder for main config
	configBuilder := config.New().AddFeeder(feeder.Env{})
	if environment == env.Development {
		configBuilder.AddFeeder(feeder.DotEnv{Path: getDotEnvPath()})
	}

	// Feed the main config struct
	if err := configBuilder.AddStruct(&cfg).Feed(); err != nil {
		return Config{}, fmt.Errorf("configBuilder.AddStruct: %w", err)
	}

	return cfg, nil
}

func MustNewConfig(environment env.Environment) Config {
	cfg, err := NewConfig(environment)
	if err != nil {
		log.Fatal("failed to load configuration: ", err)
	}
	return cfg
}

func (c *Config) DBConnectionString() string {
	if c.Db.ConnectionString != "" {
		return c.Db.ConnectionString
	}
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", c.Db.Host, c.Db.Port, c.Db.User, c.Db.Name, c.Db.Password)
}

func getDotEnvPath() string {
	exPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting dotenv path", err)
	}
	return path.Join(exPath, ".env")
}
