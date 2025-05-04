package controller

import (
	"fmt"
	"net/http"

	"github.com/Beigelman/nossas-despesas/internal/modules/group/query"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
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
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		balances, err := getGroupBalance(ctx.Context(), groupID)
		if err != nil {
			return fmt.Errorf("query.GetGroupBalance: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, GetGroupBalanceResponse{
			GroupID:  groupID,
			Balances: balances,
		}))
	}
}
