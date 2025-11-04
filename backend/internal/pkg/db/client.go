package db

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	conn          *sqlx.DB
	migrateClient *migrate.Migrate
	cfg           dbConfig
}

func NewClient(connString string, options ...Option) (*Client, error) {
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		slog.Error("sqlx.Connect", "error", err)
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	cfg := defaultConfig()
	for _, opt := range options {
		if err := opt(&cfg); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	dbName := getDBName(connString)
	cfg.name = dbName

	conn.SetMaxIdleConns(cfg.maxIdleConns)
	conn.SetMaxOpenConns(cfg.maxOpenConns)
	conn.SetConnMaxIdleTime(cfg.connMaxIdleTime)
	conn.SetConnMaxLifetime(cfg.connMaxLifeTime)

	client := &Client{
		conn:          conn,
		migrateClient: nil,
		cfg:           cfg,
	}

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
			SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' AND table_name != 'schema_migrations';
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
