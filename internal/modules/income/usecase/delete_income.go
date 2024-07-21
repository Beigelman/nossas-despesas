package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"

	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	DeleteIncomeParams struct {
		ID      income.ID
		UserID  user.ID
		GroupID group.ID
	}
	DeleteIncome func(ctx context.Context, p DeleteIncomeParams) (*income.Income, error)
)

func NewDeleteIncome(
	incomeRepo income.Repository,
) DeleteIncome {
	return func(ctx context.Context, p DeleteIncomeParams) (*income.Income, error) {
		inc, err := incomeRepo.GetByID(ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("incomeRepo.GetByID: %w", err)
		}

		// TODO: Bypass para permitir a Lu editar minhas receitas. Pensar em uma solução mais estruturante no futuro
		if p.GroupID.Value != 1 && inc.UserID != p.UserID {
			return nil, except.ForbiddenError("user mismatch")
		}

		inc.Delete()

		if err := incomeRepo.Store(ctx, inc); err != nil {
			return nil, fmt.Errorf("incomeRepo.Store: %w", err)
		}

		return inc, nil
	}
}
