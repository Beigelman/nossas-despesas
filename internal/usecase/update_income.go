package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	UpdateIncomeParams struct {
		ID        entity.IncomeID
		UserID    entity.UserID
		GroupID   entity.GroupID
		Type      *entity.IncomeType
		Amount    *int
		CreatedAt *time.Time
	}
	UpdateIncome func(ctx context.Context, p UpdateIncomeParams) (*entity.Income, error)
)

func NewUpdateIncome(
	incomeRepo repository.IncomeRepository,
) UpdateIncome {
	return func(ctx context.Context, p UpdateIncomeParams) (*entity.Income, error) {
		income, err := incomeRepo.GetByID(ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("incomeRepo.GetByID: %w", err)
		}

		// TODO: Bypass para permitir a Lu editar minhas receitas. Pensar em uma solução mais estruturante no futuro
		if p.GroupID.Value != 1 && income.UserID != p.UserID {
			return nil, except.BadRequestError("user mismatch")
		}

		income.Update(entity.UpdateIncomeParams{
			Amount:    p.Amount,
			Type:      p.Type,
			CreatedAt: p.CreatedAt,
		})

		if err := incomeRepo.Store(ctx, income); err != nil {
			return nil, fmt.Errorf("incomeRepo.Store: %w", err)
		}

		return income, nil
	}
}
