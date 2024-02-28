package handler

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type (
	DeleteIncome func(ctx *fiber.Ctx) error

	DeleteIncomeResponse struct {
		ID int `json:"id"`
	}
)

func NewDeleteIncome(deleteIncome usecase.DeleteIncome) DeleteIncome {
	return func(ctx *fiber.Ctx) error {
		userID, ok := ctx.Locals("user_id").(int)
		if !ok {
			return except.BadRequestError("invalid user id")
		}

		incomeID, err := strconv.Atoi(ctx.Params("income_id"))
		if err != nil {
			return except.BadRequestError("invalid income id")
		}

		income, err := deleteIncome(ctx.Context(), usecase.DeleteIncomeParams{
			ID:     entity.IncomeID{Value: incomeID},
			UserID: entity.UserID{Value: userID},
		})
		if err != nil {
			return fmt.Errorf("updateIncome: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, DeleteIncomeResponse{ID: income.ID.Value}),
		)
	}
}
