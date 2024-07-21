package controller

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/query"
	"net/http"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
)

type (
	GetMonthlyIncome func(ctx *fiber.Ctx) error

	GetMonthlyIncomeResponse struct {
		GroupID int                `json:"group_id"`
		Incomes []query.UserIncome `json:"incomes"`
		Total   int                `json:"total"`
		Month   time.Month         `json:"month"`
	}
)

func NewGetMonthlyIncome(getGroupMonthlyIncome query.GetMonthlyIncome) GetMonthlyIncome {
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}
		date, err := time.Parse(time.DateOnly, ctx.Query("date", ""))
		if err != nil {
			return except.BadRequestError("invalid date")
		}

		incs, err := getGroupMonthlyIncome(ctx.Context(), groupID, date)
		if err != nil {
			return fmt.Errorf("query.GetGroupMonthlyIncome: %w", err)
		}

		var totalIncome int
		for _, inc := range incs {
			totalIncome += inc.Amount
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, GetMonthlyIncomeResponse{
			GroupID: groupID,
			Incomes: incs,
			Total:   totalIncome,
			Month:   date.Month(),
		}))
	}
}
