package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Beigelman/nossas-despesas/internal/pkg/pubsub"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
)

type RecalculateExpensesSplitRatio func(ctx context.Context) error

func NewRecalculateExpensesSplitRatio(
	subscriber pubsub.Subscriber,
	recalculateExpenses usecase.RecalculateExpensesSplitRatio,
) RecalculateExpensesSplitRatio {
	return func(ctx context.Context) error {
		messages, err := subscriber.Subscribe(ctx, pubsub.IncomesTopic)
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

				if err := recalculateExpenses(ctx, usecase.RecalculateExpensesSplitRatioInput{
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
}
