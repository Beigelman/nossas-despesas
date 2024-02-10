package db

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"

	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	db            *sqlx.DB
	env           env.Environment
	name          string
	migrateClient *migrate.Migrate
	migrationPath string
	kind          string
}

func New(c *config.Config) Database {
	db, err := sqlx.Open(c.Db.Type, c.DBConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	if c.Env == env.Production {
		db.DB.SetMaxIdleConns(c.Db.MaxIdleConns)
		db.DB.SetMaxOpenConns(c.Db.MaxOpenConns)
		db.DB.SetConnMaxLifetime(c.Db.MaxLifeTime)
		db.DB.SetConnMaxIdleTime(c.Db.MaxIdleTime)
	}

	return &SQLDatabase{
		db:            db,
		env:           c.Env,
		name:          c.Db.Name,
		migrationPath: c.Db.MigrationPath,
		kind:          c.Db.Type,
	}
}

func (sql *SQLDatabase) Client() *sqlx.DB {
	return sql.db
}

func (sql *SQLDatabase) Close() error {
	return sql.db.Close()
}

func (sql *SQLDatabase) Clean(tables ...string) error {
	if len(tables) == 0 {
		rows, err := sql.db.Queryx(`
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

			if _, err = sql.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", tableName)); err != nil {
				return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
			}
		}

		return nil
	}

	for _, table := range tables {
		if _, err := sql.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", table)); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}
