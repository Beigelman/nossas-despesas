package controller

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/query"
	"net/http"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
)

type (
	GetExpensesPerSearch func(ctx *fiber.Ctx) error
)

func NewGetExpensesPerSearch(getExpensesPerSearch query.GetExpensesPerSearch) GetExpensesPerSearch {
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		search := ctx.Query("search")

		expenses, err := getExpensesPerSearch(ctx.Context(), groupID, search)
		if err != nil {
			return fmt.Errorf("query.getGroup: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, expenses))
	}
}
