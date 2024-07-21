package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"time"
)

type (
	CreateIncomeParams struct {
		Type      income.Type
		Amount    int
		UserID    user.ID
		CreatedAt *time.Time
	}
	CreateIncome func(ctx context.Context, p CreateIncomeParams) (*income.Income, error)
)

func NewCreateIncome(
	userRepo user.Repository,
	incomeRepo income.Repository,
) CreateIncome {
	return func(ctx context.Context, p CreateIncomeParams) (*income.Income, error) {
		usr, err := userRepo.GetByID(ctx, p.UserID)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if usr == nil {
			return nil, except.NotFoundError("user not found")
		}

		inc := income.New(income.Attributes{
			ID:        incomeRepo.GetNextID(),
			UserID:    usr.ID,
			Amount:    p.Amount,
			Type:      p.Type,
			CreatedAt: p.CreatedAt,
		})

		if err := incomeRepo.Store(ctx, inc); err != nil {
			return nil, fmt.Errorf("incomeRepo.Store: %w", err)
		}

		return inc, nil
	}
}
