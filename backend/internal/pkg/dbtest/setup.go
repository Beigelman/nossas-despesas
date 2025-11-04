package dbtest

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

var (
	container *db.PostgresContainer
	once      sync.Once
	dbCounter int
	mutex     sync.Mutex
)

func Setup(ctx context.Context, t *testing.T) *db.Client {
	t.Helper()

	var err error
	once.Do(func() {
		container, err = db.StartPostgres(ctx)
		require.NoError(t, err)
	})

	// Generate a unique database name for this test
	dbName := generateUniqueDBName()

	defaultDB, err := sqlx.Connect("pgx", container.ConnString())
	require.NoError(t, err)

	// Create the new test database
	err = createDatabase(defaultDB, dbName)
	require.NoError(t, err)

	// Create connection string for the new test database
	testConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		container.Host, container.Port, container.User, container.Password, dbName)

	// Create client for the test database
	testDB, err := db.NewClient(testConnString, db.WithMigrationPath("database/migrations"))
	require.NoError(t, err)

	require.NoError(t, testDB.MigrateUp())

	t.Cleanup(func() {
		require.NoError(t, testDB.Close())
		require.NoError(t, dropDatabase(defaultDB, dbName))
		require.NoError(t, defaultDB.Close())
	})

	return testDB
}

// generateUniqueDBName generates a unique database name for each test
func generateUniqueDBName() string {
	mutex.Lock()
	defer mutex.Unlock()

	dbCounter++
	return fmt.Sprintf("test_db_%d", dbCounter)
}

// createDatabase creates a new database with the given name.
// It connects to the default 'postgres' database to create the new database.
func createDatabase(conn *sqlx.DB, dbName string) error {
	if _, err := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
		return fmt.Errorf("failed to create database %s: %w", dbName, err)
	}
	return nil
}

// dropDatabase drops the database with the given name.
// It connects to the default 'postgres' database to drop the target database.
func dropDatabase(conn *sqlx.DB, dbName string) error {
	terminateQuery := fmt.Sprintf(`
	SELECT pg_terminate_backend(pid)
	FROM pg_stat_activity
	WHERE datname = '%s' AND pid <> pg_backend_pid()`, dbName)

	_, err := conn.Exec(terminateQuery)
	if err != nil {
		return fmt.Errorf("failed to terminate connections to database %s: %w", dbName, err)
	}

	if _, err := conn.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)); err != nil {
		return fmt.Errorf("failed to drop database %s: %w", dbName, err)
	}

	return nil
}
