package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"time"
)

type (
	UpdateExpenseParams struct {
		ID          entity.ExpenseID
		Name        *string
		Amount      *int
		Description *string
		CategoryID  *entity.CategoryID
		SplitRatio  *entity.SplitRatio
		PayerID     *entity.UserID
		ReceiverID  *entity.UserID
		CreatedAt   *time.Time
	}
	UpdateExpense func(ctx context.Context, p UpdateExpenseParams) (*entity.Expense, error)
)

func NewUpdateExpense(
	expenseRepo repository.ExpenseRepository,
	userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository,
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

		if err := expense.Update(entity.ExpenseUpdateParams{
			Name:        p.Name,
			Amount:      p.Amount,
			Description: p.Description,
			CategoryID:  p.CategoryID,
			SplitRatio:  p.SplitRatio,
			PayerID:     p.PayerID,
			ReceiverID:  p.ReceiverID,
			CreatedAt:   p.CreatedAt,
		}); err != nil {
			return nil, except.UnprocessableEntityError().SetInternal(fmt.Errorf("entity.Update: %w", err))
		}

		if err := expenseRepo.Store(ctx, expense); err != nil {
			return nil, fmt.Errorf("expenseRepo.Store: %w", err)
		}

		return expense, nil
	}
}
