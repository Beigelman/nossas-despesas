package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	vo "github.com/Beigelman/ludaapi/internal/domain/valueobject"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"time"
)

type (
	UpdateExpenseParams struct {
		ID           entity.ExpenseID
		Name         *string
		Amount       *int
		RefundAmount *int
		Description  *string
		CategoryID   *entity.CategoryID
		SplitType    *vo.SplitType
		PayerID      *entity.UserID
		ReceiverID   *entity.UserID
		CreatedAt    *time.Time
	}
	UpdateExpense func(ctx context.Context, p UpdateExpenseParams) (*entity.Expense, error)
)

func NewUpdateExpense(
	expenseRepo repository.ExpenseRepository,
	userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository,
	incomeRepo repository.IncomeRepository,
) UpdateExpense {
	return func(ctx context.Context, p UpdateExpenseParams) (*entity.Expense, error) {
		expense, err := expenseRepo.GetByID(ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("expenseRepo.GetByID: %w", err)
		}

		if expense == nil {
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

			if payer.GroupID == nil || expense.GroupID != *payer.GroupID {
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

			if receiver.GroupID == nil || expense.GroupID != *receiver.GroupID {
				return nil, except.UnprocessableEntityError("group mismatch")
			}
		}

		if p.CategoryID != nil {
			category, err := categoryRepo.GetByID(ctx, *p.CategoryID)
			if err != nil {
				return nil, fmt.Errorf("categoryRepo.GetByID: %w", err)
			}

			if category == nil {
				return nil, except.NotFoundError("category not found")
			}
		}

		splitRatio := expense.SplitRatio
		if p.SplitType != nil && *p.SplitType != expense.SplitRatio.Type() {
			switch *p.SplitType {
			case vo.SpliteTypes.Proportional:
				payerID := expense.PayerID
				if p.PayerID != nil {
					payerID = *p.PayerID
				}
				payerIncomes, err := incomeRepo.GetUserMonthlyIncomes(ctx, payerID, p.CreatedAt)
				if err != nil || payerIncomes == nil {
					return nil, except.UnprocessableEntityError("payer income not found").SetInternal(fmt.Errorf("incomeRepo.GetUserMonthlyIncomes: %w", err))
				}

				receiverID := expense.ReceiverID
				if p.ReceiverID != nil {
					receiverID = *p.ReceiverID
				}
				receiverIncomes, err := incomeRepo.GetUserMonthlyIncomes(ctx, receiverID, p.CreatedAt)
				if err != nil || receiverIncomes == nil {
					return nil, except.UnprocessableEntityError("receiver income not found").SetInternal(fmt.Errorf("incomeRepo.GetUserMonthlyIncomes: %w", err))
				}

				totalPayerIncome := 0
				for _, income := range payerIncomes {
					totalPayerIncome += income.Amount
				}

				totalReceiverIncome := 0
				for _, income := range receiverIncomes {
					totalReceiverIncome += income.Amount
				}

				splitRatio = vo.NewProportionalSplitRatio(totalPayerIncome, totalReceiverIncome)
			case vo.SpliteTypes.Transfer:
				splitRatio = vo.NewTransferRatio()
			default:
				splitRatio = vo.NewEqualSplitRatio()
			}
		}

		if err := expense.Update(entity.ExpenseUpdateParams{
			Name:         p.Name,
			Amount:       p.Amount,
			RefundAmount: p.RefundAmount,
			Description:  p.Description,
			CategoryID:   p.CategoryID,
			SplitRatio:   &splitRatio,
			PayerID:      p.PayerID,
			ReceiverID:   p.ReceiverID,
			CreatedAt:    p.CreatedAt,
		}); err != nil {
			return nil, except.UnprocessableEntityError().SetInternal(fmt.Errorf("expense.Update: %w", err))
		}

		if err := expenseRepo.Store(ctx, expense); err != nil {
			return nil, fmt.Errorf("expenseRepo.Store: %w", err)
		}

		return expense, nil
	}
}
