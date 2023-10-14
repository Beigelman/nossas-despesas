package db

import (
	"fmt"
	"log"

	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	Client() *sqlx.DB
	Close() error
	Clean() error
	MigrateUp(migrationPath string) error
	MigrateDown(migrationPath string) error
	NewTransactionManager() TransactionManager
}

type SQLDatabase struct {
	db   *sqlx.DB
	env  env.Environment
	name string
}

func New(c *config.Config) Database {
	db, err := sqlx.Open(c.Db.Type, c.DBConnectionURL())
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
		env:  c.Env,
		db:   db,
		name: c.Db.Name,
	}
}

func (sql *SQLDatabase) Client() *sqlx.DB {
	return sql.db
}

func (sql *SQLDatabase) Close() error {
	return sql.db.Close()
}

func (sql *SQLDatabase) Clean() error {
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

		sql_command := fmt.Sprintf("TRUNCATE TABLE %s;", tableName)
		_, err = sql.db.Exec(sql_command)
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
		}
	}
	return nil
}
