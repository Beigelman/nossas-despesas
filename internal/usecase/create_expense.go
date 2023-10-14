package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
)

type (
	CreateExpenseParams struct {
		GroupID     entity.GroupID
		Name        string
		Amount      int
		Description string
		CategoryID  entity.CategoryID
		SplitRatio  entity.SplitRatio
		PayerID     entity.UserID
		ReceiverID  entity.UserID
	}
	CreateExpense func(ctx context.Context, p CreateExpenseParams) (*entity.Expense, error)
)

func NewCreateExpense(
	expenseRepo repository.ExpenseRepository,
	userRepo repository.UserRepository,
	groupRepo repository.GroupRepository,
	categoryRepo repository.CategoryRepository,
) CreateExpense {
	return func(ctx context.Context, p CreateExpenseParams) (*entity.Expense, error) {
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

		group, err := groupRepo.GetByID(ctx, p.GroupID)
		if err != nil {
			return nil, fmt.Errorf("groupRepo.GetByID: %w", err)
		}

		if group == nil {
			return nil, except.NotFoundError("group not found")
		}

		if group.ID != *payer.GroupID || group.ID != *receiver.GroupID {
			return nil, except.UnprocessableEntityError("group mismatch")
		}

		category, err := categoryRepo.GetByID(ctx, p.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("categoryRepo.GetByID: %w", err)
		}

		if category == nil {
			return nil, except.NotFoundError("category not found")
		}

		expenseID := expenseRepo.GetNextID()

		expense, err := entity.NewExpense(entity.ExpenseParams{
			ID:          expenseID,
			Name:        p.Name,
			Amount:      p.Amount,
			Description: p.Description,
			GroupID:     p.GroupID,
			CategoryID:  p.CategoryID,
			SplitRatio:  p.SplitRatio,
			PayerID:     p.PayerID,
			ReceiverID:  p.ReceiverID,
		})
		if err != nil {
			return nil, except.UnprocessableEntityError().SetInternal(fmt.Errorf("entity.NewCategory: %w", err))
		}

		if err := expenseRepo.Store(ctx, expense); err != nil {
			return nil, fmt.Errorf("expenseRepo.Store: %w", err)
		}

		return expense, nil
	}
}
