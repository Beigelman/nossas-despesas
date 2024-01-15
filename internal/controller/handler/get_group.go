package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/query"
	"github.com/gofiber/fiber/v2"
	"net/http"
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

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[*query.Group](http.StatusOK, group))
	}
}
