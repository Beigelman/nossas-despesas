package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"log/slog"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
)

type (
	RecalculateExpensesSplitRatioInput struct {
		EventName string
		GroupID   group.ID
		Date      time.Time
	}

	RecalculateExpensesSplitRatio func(ctx context.Context, input RecalculateExpensesSplitRatioInput) error
)

func NewRecalculateExpensesSplitRatio(
	expenseRepo expense.Repository,
	incomeRepo income.Repository,
) RecalculateExpensesSplitRatio {
	return func(ctx context.Context, input RecalculateExpensesSplitRatioInput) error {
		slog.InfoContext(ctx, "Recalculating expenses split ratio", slog.Int("group", input.GroupID.Value), slog.Time("date", input.Date), slog.String("event", input.EventName))
		expenses, err := expenseRepo.GetByGroupDate(ctx, input.GroupID, input.Date)
		if err != nil {
			return fmt.Errorf("expensesRepo.GetByGroupDate: %w", err)
		}

		var proportionalExpenses []expense.Expense
		for _, expense := range expenses {
			if expense.SplitType == expense.SpliteTypes.Proportional {
				proportionalExpenses = append(proportionalExpenses, expense)
			}
		}

		if len(proportionalExpenses) == 0 {
			slog.InfoContext(ctx, "no expenses to update")
			return nil
		}

		usersIDs := []entity.UserID{proportionalExpenses[0].PayerID, proportionalExpenses[0].ReceiverID}
		usersIncomes := map[entity.UserID]int{}
		for _, userID := range usersIDs {
			incomes, err := incomeRepo.GetUserMonthlyIncomes(ctx, userID, &input.Date)
			if err != nil || incomes == nil {
				return fmt.Errorf("no incomes found for user %d", userID.Value)
			}

			totalIncome := 0
			for _, income := range incomes {
				totalIncome += income.Amount
			}

			usersIncomes[userID] = totalIncome
		}

		for i, expense := range proportionalExpenses {
			newSplitRatio := expense.NewProportionalSplitRatio(usersIncomes[expense.PayerID], usersIncomes[expense.ReceiverID])
			proportionalExpenses[i].Update(expense.ExpenseUpdateParams{SplitRatio: &newSplitRatio})
		}

		if err := expenseRepo.BulkStore(ctx, proportionalExpenses); err != nil {
			return fmt.Errorf("expense.BulkStore: %w", err)
		}

		slog.InfoContext(ctx, "Expenses split ratio recalculated successfully", slog.Int("count", len(proportionalExpenses)))

		return nil
	}
}
