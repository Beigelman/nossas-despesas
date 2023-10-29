package db

import (
	"fmt"

	"github.com/Beigelman/ludaapi/internal/pkg/env"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func (sql *SQLDatabase) MigrateUp(migrationPath string) error {
	driver, err := postgres.WithInstance(sql.db.DB, &postgres.Config{
		DatabaseName: sql.name,
	})
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}

	migrateClient, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate client: %w", err)
	}

	err = migrateClient.Up()
	if err != nil {
		return fmt.Errorf("failed to perform migration: %w", err)
	}
	return nil
}

func (sql *SQLDatabase) MigrateDown(migrationPath string) error {
	if sql.env == env.Production {
		return nil
	}

	driver, err := postgres.WithInstance(sql.db.DB, &postgres.Config{
		DatabaseName: sql.name,
	})
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}

	migrateClient, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate client: %w", err)
	}

	err = migrateClient.Down()
	if err != nil {
		return fmt.Errorf("failed to perform migration down: %w", err)
	}
	return nil
}
