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
	CreateIncomeParams struct {
		Type      entity.IncomeType
		Amount    int
		UserID    entity.UserID
		CreatedAt *time.Time
	}
	CreateIncome func(ctx context.Context, p CreateIncomeParams) (*entity.Income, error)
)

func NewCreateIncome(
	userRepo repository.UserRepository,
	incomeRepo repository.IncomeRepository,
) CreateIncome {
	return func(ctx context.Context, p CreateIncomeParams) (*entity.Income, error) {
		user, err := userRepo.GetByID(ctx, p.UserID)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if user == nil {
			return nil, except.NotFoundError("user not found")
		}

		expense := entity.NewIncome(entity.IncomeParams{
			ID:        incomeRepo.GetNextID(),
			UserID:    user.ID,
			Amount:    p.Amount,
			Type:      p.Type,
			CreatedAt: p.CreatedAt,
		})
		if err != nil {
			return nil, except.UnprocessableEntityError().SetInternal(fmt.Errorf("entity.NewIncome: %w", err))
		}

		if err := incomeRepo.Store(ctx, expense); err != nil {
			return nil, fmt.Errorf("incomeRepo.Store: %w", err)
		}

		return expense, nil
	}
}
