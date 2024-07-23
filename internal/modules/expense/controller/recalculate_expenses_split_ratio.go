package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/pubsub"
	"github.com/ThreeDotsLabs/watermill/message"
	"log/slog"
)

type RecalculateExpensesSplitRatio struct {
	subscriber          message.Subscriber
	recalculateExpenses usecase.RecalculateExpensesSplitRatio
}

func NewRecalculateExpensesSplitRatio(
	subs message.Subscriber,
	recalculateExpenses usecase.RecalculateExpensesSplitRatio,
) *RecalculateExpensesSplitRatio {
	return &RecalculateExpensesSplitRatio{
		subscriber:          subs,
		recalculateExpenses: recalculateExpenses,
	}
}

func (r RecalculateExpensesSplitRatio) Run(ctx context.Context) error {
	messages, err := r.subscriber.Subscribe(ctx, pubsub.IncomesTopic)
	if err != nil {
		return fmt.Errorf("subscriber.Subscribe: %w", err)
	}

	go func() {
		slog.InfoContext(ctx, "Listening to incomes topic...")
		for msg := range messages {
			var payload pubsub.IncomeEvent
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				msg.Nack()
				continue
			}

			if err := r.recalculateExpenses(ctx, usecase.RecalculateExpensesSplitRatioInput{
				EventName: payload.Type,
				GroupID:   payload.GroupID,
				Date:      payload.Income.CreatedAt,
			}); err != nil {
				slog.ErrorContext(ctx, "failed to recalculate expenses spit ratio", "error", err)
				msg.Nack()
				continue
			}

			msg.Ack()
		}
	}()

	return nil
}
