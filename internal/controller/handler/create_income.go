package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/infra/pubsub"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"github.com/Beigelman/nossas-despesas/internal/usecase"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type (
	CreateIncome func(ctx *fiber.Ctx) error

	CreateIncomeRequest struct {
		Type      string     `json:"type" validate:"oneof=salary benefit vacation thirteenth_salary other"`
		Amount    int        `json:"amount" validate:"required"`
		CreatedAt *time.Time `json:"created_at"`
		UserID    *int       `json:"user_id"`
	}

	CreateIncomeResponse struct {
		ID int `json:"id"`
	}
)

func NewCreateIncome(createIncome usecase.CreateIncome, publisher message.Publisher) CreateIncome {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req CreateIncomeRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		var (
			userID int
			ok     bool
		)
		if req.UserID == nil {
			userID, ok = ctx.Locals("user_id").(int)
			if !ok {
				return except.BadRequestError("invalid user id")
			}
		} else {
			userID = *req.UserID
		}

		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("missing context group_id")
		}

		income, err := createIncome(ctx.Context(), usecase.CreateIncomeParams{
			UserID:    entity.UserID{Value: userID},
			Type:      entity.IncomeType(req.Type),
			Amount:    req.Amount,
			CreatedAt: req.CreatedAt,
		})
		if err != nil {
			return fmt.Errorf("createIncome: %w", err)
		}

		event, err := json.Marshal(pubsub.IncomeEvent{
			Event: pubsub.Event{
				SentAt:  time.Now(),
				Type:    "income_created",
				UserID:  entity.UserID{Value: userID},
				GroupID: entity.GroupID{Value: groupID},
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
			api.NewResponse(http.StatusCreated, CreateIncomeResponse{ID: income.ID.Value}),
		)
	}
}
