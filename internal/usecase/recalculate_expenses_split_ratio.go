package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	vo "github.com/Beigelman/nossas-despesas/internal/domain/valueobject"
)

type (
	RecalculateExpensesSplitRatioInput struct {
		GroupID entity.GroupID
		Date    time.Time
	}

	RecalculateExpensesSplitRatio func(ctx context.Context, input RecalculateExpensesSplitRatioInput) error
)

func NewRecalculateExpensesSplitRatio(
	expenseRepo repository.ExpenseRepository,
	incomeRepo repository.IncomeRepository,
) RecalculateExpensesSplitRatio {
	return func(ctx context.Context, input RecalculateExpensesSplitRatioInput) error {
		expenses, err := expenseRepo.GetByGroupDate(ctx, input.GroupID, input.Date)
		if err != nil {
			return fmt.Errorf("expensesRepo.GetByGroupDate: %w", err)
		}

		var proportionalExpenses []entity.Expense
		for _, expense := range expenses {
			if expense.SplitRatio.Type() == vo.SpliteTypes.Proportional {
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
			newSplitRatio := vo.NewProportionalSplitRatio(usersIncomes[expense.PayerID], usersIncomes[expense.ReceiverID])
			proportionalExpenses[i].Update(entity.ExpenseUpdateParams{SplitRatio: &newSplitRatio})
		}

		if err := expenseRepo.BulkStore(ctx, proportionalExpenses); err != nil {
			return fmt.Errorf("expense.BulkStore: %w", err)
		}

		return nil
	}
}
