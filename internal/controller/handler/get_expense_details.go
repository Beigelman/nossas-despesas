package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/query"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type (
	GetExpenseDetails func(ctx *fiber.Ctx) error
)

func NewGetExpenseDetails(getGroupExpenses query.GetExpenseDetails) GetExpenseDetails {
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		expenseID, err := strconv.Atoi(ctx.Params("expense_id"))
		if err != nil {
			return except.BadRequestError("invalid expense id")
		}

		expenseDetails, err := getGroupExpenses(ctx.Context(), expenseID)
		if err != nil {
			return fmt.Errorf("query.GetExpenseDetails: %w", err)
		}

		if len(expenseDetails) == 0 {
			return except.NotFoundError("expense not found")
		}

		if expenseDetails[0].GroupID != groupID {
			return except.ForbiddenError("group mismatch")
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[[]query.ExpenseDetails](http.StatusOK, expenseDetails))
	}
}
