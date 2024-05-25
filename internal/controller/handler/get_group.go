package handler

import (
	"fmt"
	"net/http"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/query"
	"github.com/gofiber/fiber/v2"
)

type (
	GetGroup func(ctx *fiber.Ctx) error
)

func NewGetGroup(getGroup query.GetGroup) GetGroup {
	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		group, err := getGroup(ctx.Context(), groupID)
		if err != nil {
			return fmt.Errorf("query.getGroup: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, group))
	}
}
