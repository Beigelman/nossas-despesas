package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/pkg/validator"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"time"
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

func NewUpdateIncome(updateIncome usecase.UpdateIncome) UpdateIncome {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
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

		income, err := updateIncome(ctx.Context(), usecase.UpdateIncomeParams{
			ID:     entity.IncomeID{Value: incomeID},
			UserID: entity.UserID{Value: userID},
			Type: func() *entity.IncomeType {
				if req.Type == nil {
					return nil
				}
				t := entity.IncomeType(*req.Type)
				return &t
			}(),
			Amount:    req.Amount,
			CreatedAt: req.CreatedAt,
		})
		if err != nil {
			return fmt.Errorf("updateIncome: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, UpdateIncomeResponse{ID: income.ID.Value}),
		)
	}
}
