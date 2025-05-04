package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
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

		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		incomeID, err := strconv.Atoi(ctx.Params("income_id"))
		if err != nil {
			return except.BadRequestError("invalid income id")
		}

		inc, err := deleteIncome(ctx.Context(), usecase.DeleteIncomeParams{
			ID:      income.ID{Value: incomeID},
			UserID:  user.ID{Value: userID},
			GroupID: group.ID{Value: groupID},
		})
		if err != nil {
			return fmt.Errorf("updateIncome: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, DeleteIncomeResponse{ID: inc.ID.Value}),
		)
	}
}
