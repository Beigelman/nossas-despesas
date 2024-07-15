package tests

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MySqlContainerOptions struct {
	testcontainers.ContainerRequest
}

type MySqlContainer struct {
	testcontainers.Container
	Host string
	Port string
}

func DefaultMySqlContainerOptions() MySqlContainerOptions {
	return MySqlContainerOptions{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mysql:8.0.36-debian",
			ExposedPorts: []string{"3306/tcp"},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort("3306/tcp"),
			),
			Env: map[string]string{
				"MYSQL_DATABASE":      "test",
				"MYSQL_ROOT_PASSWORD": "root",
			},
		},
	}
}

func StartMySql(ctx context.Context) (*MySqlContainer, error) {
	return StartMySqlWithConfig(ctx, DefaultMySqlContainerOptions())
}

func StartMySqlWithConfig(ctx context.Context, options MySqlContainerOptions) (*MySqlContainer, error) {
	req := options.ContainerRequest
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           DisableLog{},
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &MySqlContainer{
		Container: container,
		Host:      hostIP,
		Port:      mappedPort.Port(),
	}, nil
}
