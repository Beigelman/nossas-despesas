package usecase

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	DeleteIncomeParams struct {
		ID      entity.IncomeID
		UserID  entity.UserID
		GroupID entity.GroupID
	}
	DeleteIncome func(ctx context.Context, p DeleteIncomeParams) (*entity.Income, error)
)

func NewDeleteIncome(
	incomeRepo repository.IncomeRepository,
) DeleteIncome {
	return func(ctx context.Context, p DeleteIncomeParams) (*entity.Income, error) {
		income, err := incomeRepo.GetByID(ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("incomeRepo.GetByID: %w", err)
		}

		// TODO: Bypass para permitir a Lu editar minhas receitas. Pensar em uma solução mais estruturante no futuro
		if p.GroupID.Value != 1 && income.UserID != p.UserID {
			return nil, except.ForbiddenError("user mismatch")
		}

		income.Delete()

		if err := incomeRepo.Store(ctx, income); err != nil {
			return nil, fmt.Errorf("incomeRepo.Store: %w", err)
		}

		return income, nil
	}
}
