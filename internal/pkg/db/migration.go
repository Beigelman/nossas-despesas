package db

import (
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func (sql *SQLDatabase) MigrateUp() error {
	migrateClient, err := sql.getMigrateClient()
	if err != nil {
		return fmt.Errorf("failed to get migrate client: %w", err)
	}

	err = migrateClient.Up()
	if err != nil {
		return fmt.Errorf("failed to perform migration: %w", err)
	}
	return nil
}

func (sql *SQLDatabase) MigrateDown() error {
	if sql.env == env.Production {
		return nil
	}

	migrateClient, err := sql.getMigrateClient()
	if err != nil {
		return fmt.Errorf("failed to get migrate client: %w", err)
	}

	err = migrateClient.Down()
	if err != nil {
		return fmt.Errorf("failed to perform migration down: %w", err)
	}
	return nil
}

func (sql *SQLDatabase) getMigrateClient() (*migrate.Migrate, error) {
	if sql.migrateClient != nil {
		return sql.migrateClient, nil
	}

	driver, err := postgres.WithInstance(sql.db.DB, &postgres.Config{
		DatabaseName: sql.name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}

	migrateClient, err := migrate.NewWithDatabaseInstance(sql.migrationPath, sql.kind, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate client: %w", err)
	}

	return migrateClient, nil
}
