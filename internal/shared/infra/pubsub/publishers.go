package pubsub

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	pubsubSql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/jmoiron/sqlx"
)

func NewSqlPublisher(db *sqlx.DB) (*pubsubSql.Publisher, error) {
	logger := watermill.NewSlogLogger(nil)
	publisher, err := pubsubSql.NewPublisher(
		db,
		pubsubSql.PublisherConfig{
			SchemaAdapter:        pubsubSql.DefaultPostgreSQLSchema{},
			AutoInitializeSchema: true,
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("pubsubSql.NewPublisher: %w", err)
	}

	return publisher, nil
}
