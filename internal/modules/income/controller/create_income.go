package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
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

func NewCreateIncome(createIncome usecase.CreateIncome) CreateIncome {
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

		inc, err := createIncome(ctx.Context(), usecase.CreateIncomeParams{
			UserID:    user.ID{Value: userID},
			GroupID:   group.ID{Value: groupID},
			Type:      income.Type(req.Type),
			Amount:    req.Amount,
			CreatedAt: req.CreatedAt,
		})
		if err != nil {
			return fmt.Errorf("createIncome: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, CreateIncomeResponse{ID: inc.ID.Value}),
		)
	}
}
