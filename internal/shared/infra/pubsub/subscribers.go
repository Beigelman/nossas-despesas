package pubsub

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/ThreeDotsLabs/watermill"
	pubsubSql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)
	Close() error
}

// SQL Subscriber

type SqlSubscriber struct {
	subscriber *pubsubSql.Subscriber
}

type Message = message.Message

func NewSqlSubscriber(db *db.Client) (Subscriber, error) {
	logger := watermill.NewSlogLogger(nil)

	subscriber, err := pubsubSql.NewSubscriber(db.Client(), pubsubSql.SubscriberConfig{
		SchemaAdapter:    pubsubSql.DefaultPostgreSQLSchema{},
		OffsetsAdapter:   pubsubSql.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("pubsubSql.NewSubscriber: %w", err)
	}

	return &SqlSubscriber{subscriber: subscriber}, nil
}

func (s SqlSubscriber) Subscribe(ctx context.Context, topic string) (<-chan *Message, error) {
	return s.subscriber.Subscribe(ctx, topic)
}

func (s SqlSubscriber) Close() error {
	return s.subscriber.Close()
}
