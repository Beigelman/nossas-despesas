package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	UpdateExpenseParams struct {
		ID           expense.ID
		Name         *string
		Amount       *int
		RefundAmount *int
		Description  *string
		CategoryID   *category.ID
		SplitType    *expense.SplitType
		PayerID      *user.ID
		ReceiverID   *user.ID
		CreatedAt    *time.Time
	}
	UpdateExpense func(ctx context.Context, p UpdateExpenseParams) (*expense.Expense, error)
)

func NewUpdateExpense(
	expenseRepo expense.Repository,
	userRepo user.Repository,
	categoryRepo category.Repository,
	incomeRepo income.Repository,
) UpdateExpense {
	return func(ctx context.Context, p UpdateExpenseParams) (*expense.Expense, error) {
		expns, err := expenseRepo.GetByID(ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("expenseRepo.GetByID: %w", err)
		}

		if expns == nil {
			return nil, except.NotFoundError("expense not found")
		}

		if p.PayerID != nil {
			payer, err := userRepo.GetByID(ctx, *p.PayerID)
			if err != nil {
				return nil, fmt.Errorf("userRepo.GetByID: %w", err)
			}

			if payer == nil {
				return nil, except.NotFoundError("payer not found")
			}

			if payer.GroupID == nil || expns.GroupID != *payer.GroupID {
				return nil, except.UnprocessableEntityError("group mismatch")
			}
		}

		if p.ReceiverID != nil {
			receiver, err := userRepo.GetByID(ctx, *p.ReceiverID)
			if err != nil {
				return nil, fmt.Errorf("userRepo.GetByID: %w", err)
			}

			if receiver == nil {
				return nil, except.NotFoundError("receiver not found")
			}

			if receiver.GroupID == nil || expns.GroupID != *receiver.GroupID {
				return nil, except.UnprocessableEntityError("group mismatch")
			}
		}

		if p.CategoryID != nil {
			catgry, err := categoryRepo.GetByID(ctx, *p.CategoryID)
			if err != nil {
				return nil, fmt.Errorf("categoryRepo.GetByID: %w", err)
			}

			if catgry == nil {
				return nil, except.NotFoundError("category not found")
			}
		}

		var splitRatio *expense.SplitRatio
		if p.SplitType != nil && *p.SplitType != expns.SplitType {
			switch *p.SplitType {
			case expense.SplitTypes.Proportional:
				createdAt := &expns.CreatedAt
				if p.CreatedAt != nil {
					createdAt = p.CreatedAt
				}

				payerID := expns.PayerID
				if p.PayerID != nil {
					payerID = *p.PayerID
				}

				payerIncomes, err := incomeRepo.GetUserMonthlyIncomes(ctx, payerID, createdAt)
				if err != nil || payerIncomes == nil {
					return nil, except.UnprocessableEntityError("payer income not found").SetInternal(fmt.Errorf("incomeRepo.GetUserMonthlyIncomes: %w", err))
				}

				receiverID := expns.ReceiverID
				if p.ReceiverID != nil {
					receiverID = *p.ReceiverID
				}
				receiverIncomes, err := incomeRepo.GetUserMonthlyIncomes(ctx, receiverID, createdAt)
				if err != nil || receiverIncomes == nil {
					return nil, except.UnprocessableEntityError("receiver income not found").SetInternal(fmt.Errorf("incomeRepo.GetUserMonthlyIncomes: %w", err))
				}

				totalPayerIncome := 0
				for _, incm := range payerIncomes {
					totalPayerIncome += incm.Amount
				}

				totalReceiverIncome := 0
				for _, incm := range receiverIncomes {
					totalReceiverIncome += incm.Amount
				}

				split := expense.NewProportionalSplitRatio(totalPayerIncome, totalReceiverIncome)
				splitRatio = &split
			case expense.SplitTypes.Transfer:
				split := expense.NewTransferRatio()
				splitRatio = &split
			default:
				split := expense.NewEqualSplitRatio()
				splitRatio = &split
			}
		}

		if err := expns.Update(expense.UpdateAttributes{
			Name:         p.Name,
			Amount:       p.Amount,
			RefundAmount: p.RefundAmount,
			Description:  p.Description,
			CategoryID:   p.CategoryID,
			SplitRatio:   splitRatio,
			SplitType:    p.SplitType,
			PayerID:      p.PayerID,
			ReceiverID:   p.ReceiverID,
			CreatedAt:    p.CreatedAt,
		}); err != nil {
			return nil, except.UnprocessableEntityError().SetInternal(fmt.Errorf("expense.Update: %w", err))
		}

		if err := expenseRepo.Store(ctx, expns); err != nil {
			return nil, fmt.Errorf("expenseRepo.Store: %w", err)
		}

		return expns, nil
	}
}
