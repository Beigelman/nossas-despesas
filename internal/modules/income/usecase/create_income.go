package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/pubsub"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	CreateIncomeParams struct {
		Type      income.Type
		Amount    int
		UserID    user.ID
		GroupID   group.ID
		CreatedAt *time.Time
	}

	CreateIncome func(ctx context.Context, p CreateIncomeParams) (*income.Income, error)
)

func NewCreateIncome(
	userRepo user.Repository,
	incomeRepo income.Repository,
	publisher pubsub.Publisher,
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

		event := pubsub.IncomeEvent{
			Event: pubsub.Event{
				SentAt:  time.Now(),
				Type:    "income_created",
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
