package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/usecase"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/pubsub"
	"net/http"
	"strconv"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type (
	DeleteIncome func(ctx *fiber.Ctx) error

	DeleteIncomeResponse struct {
		ID int `json:"id"`
	}
)

func NewDeleteIncome(deleteIncome usecase.DeleteIncome, publisher message.Publisher) DeleteIncome {
	return func(ctx *fiber.Ctx) error {
		userID, ok := ctx.Locals("user_id").(int)
		if !ok {
			return except.BadRequestError("invalid user id")
		}

		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		incomeID, err := strconv.Atoi(ctx.Params("income_id"))
		if err != nil {
			return except.BadRequestError("invalid income id")
		}

		income, err := deleteIncome(ctx.Context(), usecase.DeleteIncomeParams{
			ID:      income.ID{Value: incomeID},
			UserID:  entity.UserID{Value: userID},
			GroupID: group.ID{Value: groupID},
		})
		if err != nil {
			return fmt.Errorf("updateIncome: %w", err)
		}

		event, err := json.Marshal(pubsub.IncomeEvent{
			Event: pubsub.Event{
				SentAt:  time.Now(),
				Type:    "income_created",
				UserID:  entity.UserID{Value: userID},
				GroupID: group.ID{Value: groupID},
			},
			Income: *income,
		})
		if err == nil {
			if err := publisher.Publish(
				pubsub.IncomesTopic,
				message.NewMessage(uuid.NewString(), event),
			); err != nil {
				slog.ErrorContext(ctx.Context(), "failed to publish income created event", "error", err)
			}
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, DeleteIncomeResponse{ID: income.ID.Value}),
		)
	}
}
