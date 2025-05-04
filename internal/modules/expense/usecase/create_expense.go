package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"

	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	CreateExpenseParams struct {
		GroupID     group.ID
		Name        string
		Amount      int
		Description string
		CategoryID  category.ID
		SplitType   expense.SplitType
		PayerID     user.ID
		ReceiverID  user.ID
		CreatedAt   *time.Time
	}
	CreateExpense func(ctx context.Context, p CreateExpenseParams) (*expense.Expense, error)
)

func NewCreateExpense(
	expenseRepo expense.Repository,
	userRepo user.Repository,
	groupRepo group.Repository,
	categoryRepo category.Repository,
	incomeRepo income.Repository,
) CreateExpense {
	return func(ctx context.Context, p CreateExpenseParams) (*expense.Expense, error) {
		payer, err := userRepo.GetByID(ctx, p.PayerID)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if payer == nil {
			return nil, except.NotFoundError("payer not found")
		}

		receiver, err := userRepo.GetByID(ctx, p.ReceiverID)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if receiver == nil {
			return nil, except.NotFoundError("receiver not found")
		}

		grp, err := groupRepo.GetByID(ctx, p.GroupID)
		if err != nil {
			return nil, fmt.Errorf("groupRepo.GetByID: %w", err)
		}

		if grp == nil {
			return nil, except.NotFoundError("group not found")
		}

		if payer.GroupID == nil || receiver.GroupID == nil || grp.ID != *payer.GroupID || grp.ID != *receiver.GroupID {
			return nil, except.UnprocessableEntityError("group mismatch")
		}

		cat, err := categoryRepo.GetByID(ctx, p.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("categoryRepo.GetByID: %w", err)
		}

		if cat == nil {
			return nil, except.NotFoundError("category not found")
		}

		var splitRatio expense.SplitRatio
		switch p.SplitType {
		case expense.SplitTypes.Proportional:
			payerIncomes, err := incomeRepo.GetUserMonthlyIncomes(ctx, payer.ID, p.CreatedAt)
			if err != nil {
				return nil, fmt.Errorf("incomeRepo.GetUserMonthlyIncomes: %w", err)
			}

			receiverIncomes, err := incomeRepo.GetUserMonthlyIncomes(ctx, receiver.ID, p.CreatedAt)
			if err != nil {
				return nil, fmt.Errorf("incomeRepo.GetUserMonthlyIncomes: %w", err)
			}

			if receiverIncomes == nil || payerIncomes == nil {
				splitRatio = expense.NewEqualSplitRatio()
			} else {
				totalPayerIncome := 0
				for _, incm := range payerIncomes {
					totalPayerIncome += incm.Amount
				}

				totalReceiverIncome := 0
				for _, incm := range receiverIncomes {
					totalReceiverIncome += incm.Amount
				}

				splitRatio = expense.NewProportionalSplitRatio(totalPayerIncome, totalReceiverIncome)
			}
		case expense.SplitTypes.Transfer:
			splitRatio = expense.NewTransferRatio()
		default:
			splitRatio = expense.NewEqualSplitRatio()
		}

		newExpense, err := expense.New(expense.Attributes{
			ID:          expenseRepo.GetNextID(),
			Name:        p.Name,
			Amount:      p.Amount,
			Description: p.Description,
			GroupID:     p.GroupID,
			CategoryID:  p.CategoryID,
			SplitRatio:  splitRatio,
			SplitType:   p.SplitType,
			PayerID:     p.PayerID,
			ReceiverID:  p.ReceiverID,
			CreatedAt:   p.CreatedAt,
		})
		if err != nil {
			return nil, except.UnprocessableEntityError().SetInternal(fmt.Errorf("entity.New: %w", err))
		}

		if err := expenseRepo.Store(ctx, newExpense); err != nil {
			return nil, fmt.Errorf("expenseRepo.Store: %w", err)
		}

		return newExpense, nil
	}
}
