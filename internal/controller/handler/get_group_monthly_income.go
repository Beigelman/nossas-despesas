package handler

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/query"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type (
	GetGroupMonthlyIncome func(ctx *fiber.Ctx) error

	GetGroupMonthlyIncomeResponse struct {
		GroupID int                `json:"group_id"`
		Incomes []query.UserIncome `json:"incomes"`
		Total   int                `json:"total"`
		Month   time.Month         `json:"month"`
	}
)

func NewGetGroupMonthlyIncome(getGroupMonthlyIncome query.GetGroupMonthlyIncome) GetGroupMonthlyIncome {
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}
		date, err := time.Parse(time.DateOnly, ctx.Query("date", ""))
		if err != nil {
			return except.BadRequestError("invalid date")
		}

		incomes, err := getGroupMonthlyIncome(ctx.Context(), groupID, date)
		if err != nil {
			return fmt.Errorf("query.GetGroupMonthlyIncome: %w", err)
		}

		var totalIncome int
		for _, income := range incomes {
			totalIncome += income.Amount
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[GetGroupMonthlyIncomeResponse](http.StatusOK, GetGroupMonthlyIncomeResponse{
			GroupID: groupID,
			Incomes: incomes,
			Total:   totalIncome,
			Month:   date.Month(),
		}))
	}
}
