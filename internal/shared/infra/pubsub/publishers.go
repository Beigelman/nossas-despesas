package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/ThreeDotsLabs/watermill"
	pubsubSql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, event any) error
	Close() error
}

// SQL Publisher

type SqlPublisher struct {
	publisher *pubsubSql.Publisher
}

func NewSqlPublisher(db *db.Client) (Publisher, error) {
	logger := watermill.NewSlogLogger(nil)
	publisher, err := pubsubSql.NewPublisher(
		db.Client(),
		pubsubSql.PublisherConfig{
			SchemaAdapter:        pubsubSql.DefaultPostgreSQLSchema{},
			AutoInitializeSchema: true,
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("pubsubSql.NewPublisher: %w", err)
	}

	return &SqlPublisher{publisher: publisher}, nil
}

func (p SqlPublisher) Publish(ctx context.Context, topic string, event any) error {
	messageData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event data: %w", err)
	}

	if err := p.publisher.Publish(topic, message.NewMessage(uuid.NewString(), messageData)); err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	return nil
}

func (p SqlPublisher) Close() error {
	return p.publisher.Close()
}
