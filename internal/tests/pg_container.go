package tests

import (
	"context"

	"github.com/docker/go-connections/nat"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainerOptions struct {
	testcontainers.ContainerRequest
}

type PostgresContainer struct {
	testcontainers.Container
	Host string
	Port string
}

func DefaultPostgresContainerOptions() PostgresContainerOptions {
	return PostgresContainerOptions{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort(nat.Port("5432/tcp")),
			),
			Env: map[string]string{
				"POSTGRES_USER":     "root",
				"POSTGRES_PASSWORD": "root",
				"POSTGRES_DB":       "test",
			},
		},
	}
}

func StartPostgres(ctx context.Context) (*PostgresContainer, error) {
	return StartPostgresWithConfig(ctx, DefaultPostgresContainerOptions())
}

type DisableLog struct{}

func (l DisableLog) Printf(format string, a ...interface{}) {}

func StartPostgresWithConfig(ctx context.Context, options PostgresContainerOptions) (*PostgresContainer, error) {
	req := options.ContainerRequest
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           DisableLog{},
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432/tcp")
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
	}, nil
}
