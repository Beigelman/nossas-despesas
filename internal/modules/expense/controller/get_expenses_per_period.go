package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
)

type GetExpensesPerPeriod func(ctx *fiber.Ctx) error

type GetExpensesPerPeriodReq struct {
	Aggregate string    `query:"aggregate"`
	StartDate time.Time `query:"start_date"`
	EndDate   time.Time `query:"end_date"`
}

func NewGetExpensesPerPeriod(getExpensesPerCategory postgres.GetExpensesPerPeriod) GetExpensesPerPeriod {
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		var params GetExpensesPerPeriodReq
		if err := ctx.QueryParser(&params); err != nil {
			return except.BadRequestError().SetInternal(err)
		}

		expensesPerPeriod, err := getExpensesPerCategory(ctx.Context(), postgres.GetExpensesPerPeriodInput{
			GroupID:   groupID,
			Aggregate: params.Aggregate,
			StartDate: params.StartDate,
			EndDate:   params.EndDate,
		})
		if err != nil {
			return fmt.Errorf("query.GetExpensesPerPeriod: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, expensesPerPeriod))
	}
}
