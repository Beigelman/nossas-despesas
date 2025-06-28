package db

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ContainerConfig represents the configuration for a Postgres container.
// It contains the image, database name, user, password, port and logger disabled option.
type ContainerConfig struct {
	Image          string
	DBName         string
	User           string
	Password       string
	Port           string
	LoggerDisabled bool
}

// PostgresContainer represents a Postgres container.
// It contains the container, the host, port, user, password and database name.

type PostgresContainer struct {
	testcontainers.Container
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// ConnString returns the connection string to the database.
func (c PostgresContainer) ConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.DBName)
}

// StartPostgres starts a new Postgres container.
// It returns a PostgresContainer or an error if the container fails to start.
// Example:
//
//	container, err := StartPostgres(ctx)
//	if err != nil {
//		log.Fatalf("failed to start postgres container: %v", err)
//	}
func StartPostgres(ctx context.Context) (*PostgresContainer, error) {
	return StartPostgresWithConfig(ctx, ContainerConfig{})
}

// StartPostgresWithConfig starts a new Postgres container with a custom configuration.
// It returns a PostgresContainer or an error if the container fails to start.
// Example:
//
//	container, err := StartPostgresWithConfig(ctx, ContainerConfig{
//		Image: "postgres:16-alpine",
//	})
func StartPostgresWithConfig(ctx context.Context, cfg ContainerConfig) (*PostgresContainer, error) {
	var logger log.Logger
	if cfg.LoggerDisabled {
		logger = noopLogger{}
	}

	req := configToContainerRequest(&cfg)
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           logger,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(req.ExposedPorts[0]))
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		Container: container,
		Host:      hostIP,
		Port:      mappedPort.Port(),
		Password:  cfg.Password,
		User:      cfg.User,
		DBName:    cfg.DBName,
	}, nil
}

func configToContainerRequest(cfg *ContainerConfig) testcontainers.ContainerRequest {
	if cfg.Image == "" {
		cfg.Image = "postgres:16-alpine"
	}

	if cfg.Port == "" {
		cfg.Port = fmt.Sprintf("%s/tcp", "5432")
	} else {
		cfg.Port = fmt.Sprintf("%s/tcp", cfg.Port)
	}

	if cfg.User == "" {
		cfg.User = "root"
	}

	if cfg.Password == "" {
		cfg.Password = "root"
	}

	if cfg.DBName == "" {
		cfg.DBName = "default"
	}

	return testcontainers.ContainerRequest{
		Image:        cfg.Image,
		ExposedPorts: []string{cfg.Port},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForSQL(nat.Port(cfg.Port), "postgres", func(host string, port nat.Port) string {
				return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, host, port.Port(), cfg.DBName)
			}),
		),
		Env: map[string]string{
			"POSTGRES_USER":     cfg.User,
			"POSTGRES_PASSWORD": cfg.Password,
			"POSTGRES_DB":       cfg.DBName,
		},
	}
}

type noopLogger struct{}

func (l noopLogger) Printf(_ string, _ ...interface{}) {}
