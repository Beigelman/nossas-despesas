package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/pubsub"
)

type CreateExpenseFromScheduled func(ctx context.Context) error

func NewCreateExpenseFromScheduled(
	subscriber pubsub.Subscriber,
	createExpense usecase.CreateExpense,
) CreateExpenseFromScheduled {
	return func(ctx context.Context) error {
		messages, err := subscriber.Subscribe(ctx, pubsub.ExpensesTopic)
		if err != nil {
			return fmt.Errorf("subscriber.Subscribe: %w", err)
		}

		go func() {
			slog.InfoContext(ctx, "Listening to expenses topic...")
			for msg := range messages {
				var payload pubsub.ExpenseEvent
				if err := json.Unmarshal(msg.Payload, &payload); err != nil {
					msg.Nack()
					continue
				}

				if _, err := createExpense(ctx, usecase.CreateExpenseParams{
					GroupID:     payload.GroupID,
					Name:        payload.Expense.Name,
					Amount:      payload.Expense.Amount,
					Description: payload.Expense.Description,
					CategoryID:  payload.Expense.CategoryID,
					SplitType:   payload.Expense.SplitType,
					PayerID:     payload.Expense.PayerID,
					ReceiverID:  payload.Expense.ReceiverID,
					CreatedAt:   &payload.Expense.CreatedAt,
				}); err != nil {
					slog.ErrorContext(ctx, "failed to create expense from scheduled", "error", err)
					msg.Nack()
					continue
				}

				msg.Ack()
			}
		}()

		return nil
	}
}
