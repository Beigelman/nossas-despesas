package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
)

type (
	DeleteIncomeParams struct {
		ID     entity.IncomeID
		UserID entity.UserID
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

		if income.UserID != p.UserID {
			return nil, except.BadRequestError("user mismatch")
		}

		income.Delete()

		if err := incomeRepo.Store(ctx, income); err != nil {
			return nil, fmt.Errorf("incomeRepo.Store: %w", err)
		}

		return income, nil
	}
}
