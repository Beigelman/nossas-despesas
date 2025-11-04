package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
)

type Option func(c *dbConfig) error

type dbConfig struct {
	env             env.Environment
	name            string
	migrationPath   string
	connMaxIdleTime time.Duration
	connMaxLifeTime time.Duration
	maxIdleConns    int
	maxOpenConns    int
}

func defaultConfig() dbConfig {
	return dbConfig{
		env:             env.Development,
		name:            "test",
		connMaxIdleTime: DefaultConnMaxIdleTime,
		connMaxLifeTime: DefaultConnMaxLifeTime,
		maxIdleConns:    DefaultMaxIdleConns,
		maxOpenConns:    DefaultMaxOpenConn,
	}
}

const (
	DefaultMaxIdleConns    = 1
	DefaultMaxOpenConn     = 3
	DefaultConnMaxLifeTime = 1 * time.Minute
	DefaultConnMaxIdleTime = 5 * time.Minute
)

// WithMaxIdleConns sets the maximum number of idle connections to the database.
// It returns an Option that can be used to configure the client.
//
// Example:
//
//	client, err := NewClient("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", sqldb.WithMaxIdleConns(10))
//	if err != nil {
//		log.Fatalf("failed to create client: %v", err)
//	}
func WithMaxIdleConns(max int) Option {
	return func(c *dbConfig) error {
		c.maxIdleConns = max
		return nil
	}
}

// WithMaxOpenConns sets the maximum number of open connections to the database.
// It returns an Option that can be used to configure the client.
//
// Example:
//
//	client, err := NewClient("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", sqldb.WithMaxOpenConns(10))
//	if err != nil {
//		log.Fatalf("failed to create client: %v", err)
//	}
func WithMaxOpenConns(max int) Option {
	return func(c *dbConfig) error {
		c.maxOpenConns = max
		return nil
	}
}

// WithConnMaxLifeTime sets the maximum lifetime of a connection to the database.
// It returns an Option that can be used to configure the client.
//
// Example:
//
//	client, err := NewClient("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", sqldb.WithConnMaxLifeTime(10 * time.Second))
//	if err != nil {
//		log.Fatalf("failed to create client: %v", err)
//	}
func WithConnMaxLifeTime(max time.Duration) Option {
	return func(c *dbConfig) error {
		c.connMaxLifeTime = max
		return nil
	}
}

// WithConnMaxIdleTime sets the maximum idle time of a connection to the database.
// It returns an Option that can be used to configure the client.
//
// Example:
//
//	client, err := NewClient("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", sqldb.WithConnMaxIdleTime(10 * time.Second))
//	if err != nil {
//		log.Fatalf("failed to create client: %v", err)
//	}
func WithConnMaxIdleTime(max time.Duration) Option {
	return func(c *dbConfig) error {
		c.connMaxIdleTime = max
		return nil
	}
}

// WithMigrationPath sets the path to the migrations files.
// It returns an Option that can be used to configure the client usually used in tests environment.
// For production environment, it's recommended to use run migrations with a CI/CD pipeline.
//
// Example:
//
//	client, err := NewClient("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", sqldb.WithMigrationPath("database/migrations"))
//	if err != nil {
//		log.Fatalf("failed to create client: %v", err)
//	}
func WithMigrationPath(path string) Option {
	return func(c *dbConfig) error {
		wrkDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("os.Getwd: %w", err)
		}

		rootDir, err := findGoModRoot(wrkDir)
		if err != nil {
			return fmt.Errorf("find root dir: %w", err)
		}

		c.migrationPath = fmt.Sprintf("file://%s/%s", rootDir, path)

		return nil
	}
}

func findGoModRoot(startDir string) (string, error) {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil // Found go.mod, return the directory
		} else if os.IsNotExist(err) {
			parentDir := filepath.Dir(dir)
			if parentDir == dir {
				return "", fmt.Errorf("go.mod not found")
			}
			dir = parentDir // Move up one directory
		} else {
			return "", err
		}
	}
}
