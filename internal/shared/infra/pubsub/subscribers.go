package pubsub

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	pubsubSql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/jmoiron/sqlx"
)

func NewSqlSubscriber(db *sqlx.DB) (*pubsubSql.Subscriber, error) {
	logger := watermill.NewSlogLogger(nil)

	subscriber, err := pubsubSql.NewSubscriber(db, pubsubSql.SubscriberConfig{
		SchemaAdapter:    pubsubSql.DefaultPostgreSQLSchema{},
		OffsetsAdapter:   pubsubSql.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("pubsubSql.NewSubscriber: %w", err)
	}

	return subscriber, nil
}
