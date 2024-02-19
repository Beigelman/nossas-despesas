package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/query"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type GetIncomesPerPeriod func(ctx *fiber.Ctx) error

type GetIncomesPerPeriodReq struct {
	StartDate time.Time `query:"start_date"`
	EndDate   time.Time `query:"end_date"`
}

func NewGetIncomesPerPeriod(getIncomesPerMonth query.GetIncomesPerPeriod) GetIncomesPerPeriod {
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		var params GetIncomesPerPeriodReq
		if err := ctx.QueryParser(&params); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		incomesPerMonth, err := getIncomesPerMonth(ctx.Context(), query.GetIncomesPerPeriodInput{
			GroupID:   groupID,
			StartDate: params.StartDate,
			EndDate:   params.EndDate,
		})
		if err != nil {
			return fmt.Errorf("query.GetIncomesPerPeriod: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[[]query.IncomesPerPeriod](http.StatusOK, incomesPerMonth))
	}
}
