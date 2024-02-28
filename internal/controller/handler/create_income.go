package handler

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"github.com/Beigelman/nossas-despesas/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type (
	CreateIncome func(ctx *fiber.Ctx) error

	CreateIncomeRequest struct {
		Type      string     `json:"type" validate:"oneof=salary benefit vacation thirteenth_salary other"`
		Amount    int        `json:"amount" validate:"required"`
		CreatedAt *time.Time `json:"created_at"`
	}

	CreateIncomeResponse struct {
		ID int `json:"id"`
	}
)

func NewCreateIncome(createIncome usecase.CreateIncome) CreateIncome {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		userID, ok := ctx.Locals("user_id").(int)
		if !ok {
			return except.BadRequestError("invalid user id")
		}

		var req CreateIncomeRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
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

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, CreateIncomeResponse{ID: income.ID.Value}),
		)
	}
}
