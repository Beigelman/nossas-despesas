package db

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	conn            *sqlx.DB
	env             env.Environment
	name            string
	migrateClient   *migrate.Migrate
	migrationPath   string
	maxIdleConns    int
	maxOpenConns    int
	connMaxIdleTime time.Duration
	connMaxLifeTime time.Duration
}

func defaultClient() *Client {
	return &Client{
		maxIdleConns:    DefaultMaxIdleConns,
		maxOpenConns:    DefaultMaxOpenConn,
		connMaxIdleTime: DefaultConnMaxIdleTime,
		connMaxLifeTime: DefaultConnMaxLifeTime,
	}
}

func NewClient(connString string, options ...Option) (*Client, error) {
	client := defaultClient()
	for _, opt := range options {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	conn.SetMaxIdleConns(client.maxIdleConns)
	conn.SetMaxOpenConns(client.maxOpenConns)
	conn.SetConnMaxIdleTime(client.connMaxIdleTime)
	conn.SetConnMaxLifetime(client.connMaxLifeTime)

	client.conn = conn
	client.name = getDBName(connString)

	return client, nil
}

func (sql *Client) Conn() *sqlx.DB {
	return sql.conn
}

func (sql *Client) Close() error {
	if sql.conn == nil {
		return nil
	}

	return sql.conn.Close()
}

func getDBName(connectionString string) string {
	re := regexp.MustCompile(`dbname=([^\s]+)`)
	matches := re.FindStringSubmatch(connectionString)
	if len(matches) > 1 {
		return matches[1]
	}

	return "default"
}

func (sql *Client) Clean(tables ...string) error {
	if len(tables) == 0 {
		rows, err := sql.conn.Queryx(`
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

			if _, err = sql.conn.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", tableName)); err != nil {
				return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
			}
		}

		return nil
	}

	for _, table := range tables {
		if _, err := sql.conn.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", table)); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}
