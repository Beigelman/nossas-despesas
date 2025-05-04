package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/pubsub"
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
	userRepo user.Repository,
	publisher pubsub.Publisher,
) UpdateIncome {
	return func(ctx context.Context, p UpdateIncomeParams) (*income.Income, error) {
		inc, err := incomeRepo.GetByID(ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("incomeRepo.GetByID: %w", err)
		}

		usr, err := userRepo.GetByID(ctx, p.UserID)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if usr == nil {
			return nil, except.NotFoundError("user not found")
		}

		if !usr.HasFlag(user.EDIT_PARTNER_INCOME) && inc.UserID != p.UserID {
			return nil, except.ForbiddenError("user mismatch")
		}

		inc.Update(income.UpdateAttributes{
			Amount:    p.Amount,
			Type:      p.Type,
			CreatedAt: p.CreatedAt,
		})

		if err := incomeRepo.Store(ctx, inc); err != nil {
			return nil, fmt.Errorf("incomeRepo.Store: %w", err)
		}

		event := pubsub.IncomeEvent{
			Event: pubsub.Event{
				SentAt:  time.Now(),
				Type:    "income_updated",
				UserID:  p.UserID,
				GroupID: p.GroupID,
			},
			Income: *inc,
		}
		if err := publisher.Publish(ctx, pubsub.IncomesTopic, event); err != nil {
			slog.ErrorContext(ctx, "failed to publish income created event", "error", err)
		}

		return inc, nil
	}
}
