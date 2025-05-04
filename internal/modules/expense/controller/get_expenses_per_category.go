package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/query"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
)

type GetExpensesPerCategory func(ctx *fiber.Ctx) error

type GetExpensesPerCategoryReq struct {
	StartDate time.Time `query:"start_date"`
	EndDate   time.Time `query:"end_date"`
}

func NewGetExpensesPerCategory(getExpensesPerCategory query.GetExpensesPerCategory) GetExpensesPerCategory {
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		var params GetExpensesPerCategoryReq
		if err := ctx.QueryParser(&params); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		expensesPerCategory, err := getExpensesPerCategory(ctx.Context(), query.GetExpensesPerCategoryInput{
			GroupID:   groupID,
			StartDate: params.StartDate,
			EndDate:   params.EndDate,
		})
		if err != nil {
			return fmt.Errorf("query.getExpensesPerCategory: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, expensesPerCategory))
	}
}
