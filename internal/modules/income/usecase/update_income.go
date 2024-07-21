package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	UpdateIncomeParams struct {
		ID        income.ID
		UserID    user.ID
		GroupID   group.ID
		Type      *income.Type
		Amount    *int
		CreatedAt *time.Time
	}
	UpdateIncome func(ctx context.Context, p UpdateIncomeParams) (*income.Income, error)
)

func NewUpdateIncome(
	incomeRepo income.Repository,
) UpdateIncome {
	return func(ctx context.Context, p UpdateIncomeParams) (*income.Income, error) {
		inc, err := incomeRepo.GetByID(ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("incomeRepo.GetByID: %w", err)
		}

		// TODO: Bypass para permitir a Lu editar minhas receitas. Pensar em uma solução mais estruturante no futuro
		if p.GroupID.Value != 1 && inc.UserID != p.UserID {
			return nil, except.BadRequestError("user mismatch")
		}

		inc.Update(income.UpdateAttributes{
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
