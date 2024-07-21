package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/shared/infra/pubsub"
	"net/http"
	"strconv"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type (
	UpdateIncome func(ctx *fiber.Ctx) error

	UpdateIncomeRequest struct {
		Type      *string    `json:"type" validate:"omitempty,oneof=salary benefit vacation thirteenth_salary other"`
		Amount    *int       `json:"amount" validate:"omitempty,gt=0"`
		CreatedAt *time.Time `json:"created_at" validate:"omitempty"`
	}

	UpdateIncomeResponse struct {
		ID int `json:"id"`
	}
)

func NewUpdateIncome(updateIncome usecase.UpdateIncome, publisher message.Publisher) UpdateIncome {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		userID, ok := ctx.Locals("user_id").(int)
		if !ok {
			return except.BadRequestError("invalid user id")
		}

		incomeID, err := strconv.Atoi(ctx.Params("income_id"))
		if err != nil {
			return except.BadRequestError("invalid income id")
		}

		var req UpdateIncomeRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		inc, err := updateIncome(ctx.Context(), usecase.UpdateIncomeParams{
			ID:      income.ID{Value: incomeID},
			UserID:  user.ID{Value: userID},
			GroupID: group.ID{Value: groupID},
			Type: func() *income.Type {
				if req.Type == nil {
					return nil
				}
				t := income.Type(*req.Type)
				return &t
			}(),
			Amount:    req.Amount,
			CreatedAt: req.CreatedAt,
		})
		if err != nil {
			return fmt.Errorf("updateIncome: %w", err)
		}

		event, err := json.Marshal(pubsub.IncomeEvent{
			Event: pubsub.Event{
				SentAt:  time.Now(),
				Type:    "income_updated",
				UserID:  user.ID{Value: userID},
				GroupID: group.ID{Value: groupID},
			},
			Income: *inc,
		})
		if err == nil {
			if err := publisher.Publish(
				pubsub.IncomesTopic,
				message.NewMessage(uuid.NewString(), event),
			); err != nil {
				slog.ErrorContext(ctx.Context(), "failed to publish income updated event")
			}
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, UpdateIncomeResponse{ID: inc.ID.Value}),
		)
	}
}
