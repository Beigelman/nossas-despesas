package db

import (
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/pkg/config"

	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	Client() *sqlx.DB
	Close() error
	Clean(tables ...string) error
	MigrateUp() error
	MigrateDown() error
	NewTransactionManager() TransactionManager
}

type SQLDatabase struct {
	*sqlx.DB
	env           env.Environment
	name          string
	migrateClient *migrate.Migrate
	migrationPath string
	kind          string
}

func New(c *config.Config) (Database, error) {
	db, err := sqlx.Connect("pgx", c.DBConnectionString())
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	if c.Env == env.Production {
		db.SetMaxIdleConns(c.Db.MaxIdleConns)
		db.SetMaxOpenConns(c.Db.MaxOpenConns)
		db.SetConnMaxLifetime(c.Db.MaxLifeTime)
		db.SetConnMaxIdleTime(c.Db.MaxIdleTime)
	}

	return &SQLDatabase{
		DB:            db,
		env:           c.Env,
		name:          c.Db.Name,
		migrationPath: c.Db.MigrationPath,
	}, nil
}

func (sql *SQLDatabase) Client() *sqlx.DB {
	return sql.DB
}

func (sql *SQLDatabase) Clean(tables ...string) error {
	if len(tables) == 0 {
		rows, err := sql.Queryx(`
			SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' and table_type = 'BASE TABLE' and table_name != 'schema_migrations';
		`)
		if err != nil {
			return fmt.Errorf("failed to get DB tables: %w", err)
		}

		var tableName string
		for rows.Next() {
			err := rows.Scan(&tableName)
			if err != nil {
				return fmt.Errorf("failed to scan %s: %w", tableName, err)
			}

			if _, err = sql.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", tableName)); err != nil {
				return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
			}
		}

		return nil
	}

	for _, table := range tables {
		if _, err := sql.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", table)); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}
