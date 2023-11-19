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
	GetGroup func(ctx *fiber.Ctx) error
)

func NewGetGroup(getGroup query.GetGroup) GetGroup {
	return func(ctx *fiber.Ctx) error {
		groupID, err := strconv.Atoi(ctx.Params("group_id"))
		if err != nil {
			return except.BadRequestError("invalid group id").SetInternal(err)
		}

		group, err := getGroup(ctx.Context(), groupID)
		if err != nil {
			return fmt.Errorf("query.getGroup: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[*query.Group](http.StatusOK, group))
	}
}
