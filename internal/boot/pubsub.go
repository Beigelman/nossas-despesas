package boot

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Beigelman/nossas-despesas/internal/infra/pubsub"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/usecase"
	"github.com/ThreeDotsLabs/watermill/message"
)

var PubSubModule = eon.NewModule("Pubsub", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, func(db db.Database) (message.Publisher, error) {
		return pubsub.NewSqlPublisher(db.Client())
	})

	di.Provide(c, func(db db.Database) (message.Subscriber, error) {
		return pubsub.NewSqlSubiscriber(db.Client())
	})

	lc.OnRunning(eon.HookOrders.APPEND, func() error {
		return di.Call(c, func(
			subs message.Subscriber,
			recalculateExpenses usecase.RecalculateExpensesSplitRatio,
		) error {
			messages, err := subs.Subscribe(ctx, pubsub.IncomesTopic)
			if err != nil {
				return fmt.Errorf("subscriber.Subscribe: %w", err)
			}

			go func() {
				slog.InfoContext(ctx, "Listeninig to incomes topic...")
				for msg := range messages {
					var payload pubsub.IncomeEvent
					if err := json.Unmarshal(msg.Payload, &payload); err != nil {
						msg.Nack()
						continue
					}

					if err := recalculateExpenses(ctx, usecase.RecalculateExpensesSplitRatioInput{
						GroupID: payload.GroupID,
						Date:    payload.Income.CreatedAt,
					}); err != nil {
						msg.Nack()
						continue
					}

					msg.Ack()
				}
			}()

			return nil
		})
	})

	lc.OnDisposing(eon.HookOrders.APPEND, func() error {
		if publisher := di.Resolve[message.Publisher](c); publisher != nil {
			slog.InfoContext(ctx, "Closing publisher connection")
			return publisher.Close()
		}
		return nil
	})

	lc.OnDisposing(eon.HookOrders.APPEND, func() error {
		if subscriber := di.Resolve[message.Subscriber](c); subscriber != nil {
			slog.InfoContext(ctx, "Closing subscriber connection")
			return subscriber.Close()
		}
		return nil
	})
})
