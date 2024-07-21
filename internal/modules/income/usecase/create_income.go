package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"time"
)

type (
	CreateIncomeParams struct {
		Type      income.Type
		Amount    int
		UserID    entity.UserID
		CreatedAt *time.Time
	}
	CreateIncome func(ctx context.Context, p CreateIncomeParams) (*income.Income, error)
)

func NewCreateIncome(
	userRepo repository.UserRepository,
	incomeRepo income.Repository,
) CreateIncome {
	return func(ctx context.Context, p CreateIncomeParams) (*income.Income, error) {
		user, err := userRepo.GetByID(ctx, p.UserID)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if user == nil {
			return nil, except.NotFoundError("user not found")
		}

		inc := income.New(income.Attributes{
			ID:        incomeRepo.GetNextID(),
			UserID:    user.ID,
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
