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
	GetGroupBalance func(ctx *fiber.Ctx) error

	GetGroupBalanceResponse struct {
		GroupID  int                 `json:"group_id"`
		Balances []query.UserBalance `json:"balances,omitempty"`
	}
)

func NewGetGroupBalance(getGroupBalance query.GetGroupBalance) GetGroupBalance {
	return func(ctx *fiber.Ctx) error {
		groupID, err := strconv.Atoi(ctx.Params("group_id"))
		if err != nil {
			return except.BadRequestError("invalid group id").SetInternal(err)
		}

		balances, err := getGroupBalance(ctx.Context(), groupID)
		if err != nil {
			return fmt.Errorf("query.GetGroupBalance: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[GetGroupBalanceResponse](http.StatusOK, GetGroupBalanceResponse{
			GroupID:  groupID,
			Balances: balances,
		}))
	}
}
