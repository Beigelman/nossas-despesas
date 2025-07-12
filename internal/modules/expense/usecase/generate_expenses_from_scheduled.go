package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/pkg/pubsub"
)

type GenerateExpensesFromScheduledUseCase func(ctx context.Context) (expensesCreated int, err error)

func NewGenerateExpensesFromScheduledUseCase(
	scheduledExpenseRepo expense.ScheduledExpenseRepository,
	publisher pubsub.Publisher,
) GenerateExpensesFromScheduledUseCase {
	return func(ctx context.Context) (expensesCreated int, err error) {
		activeScheduledExpenses, err := scheduledExpenseRepo.GetActiveScheduledExpenses(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to get active scheduled expenses: %w", err)
		}

		var (
			expenses                 []expense.Expense
			scheduledExpensesCreated []expense.ScheduledExpense
		)
		for _, scheduledExpense := range activeScheduledExpenses {
			if !scheduledExpense.ShouldGenerateExpense() {
				continue
			}

			exp, err := scheduledExpense.ToExpense()
			if err != nil {
				return 0, fmt.Errorf("failed to convert scheduled expense to expense: %w", err)
			}
			expenses = append(expenses, *exp)

			scheduledExpense.UpdateLastGeneratedAt()
			scheduledExpensesCreated = append(scheduledExpensesCreated, scheduledExpense)
		}

		if len(scheduledExpensesCreated) > 0 {
			if err := scheduledExpenseRepo.BulkStore(ctx, scheduledExpensesCreated); err != nil {
				return 0, fmt.Errorf("failed to store scheduled expenses: %w", err)
			}
		}

		for _, exp := range expenses {
			if err := publisher.Publish(ctx, pubsub.ExpensesTopic, pubsub.ExpenseEvent{
				Event: pubsub.Event{
					Type:    "expense.created",
					GroupID: exp.GroupID,
					UserID:  exp.PayerID,
					SentAt:  time.Now(),
				},
				Expense: exp,
			}); err != nil {
				return 0, fmt.Errorf("failed to publish expense created event: %w", err)
			}
		}

		return len(expenses), nil
	}
}
