package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	GetExpenseDetails func(ctx *fiber.Ctx) error
)

func NewGetExpenseDetails(getExpenseDetails postgres.GetExpenseDetails) GetExpenseDetails {
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		expenseID, err := strconv.Atoi(ctx.Params("expense_id"))
		if err != nil {
			return except.BadRequestError("invalid expense id")
		}

		expenseDetails, err := getExpenseDetails(ctx.Context(), expenseID)
		if err != nil {
			return fmt.Errorf("query.GetExpenseDetails: %w", err)
		}

		if len(expenseDetails) == 0 {
			return except.NotFoundError("expense not found")
		}

		if expenseDetails[0].GroupID != groupID {
			return except.ForbiddenError("group mismatch")
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, expenseDetails))
	}
}
