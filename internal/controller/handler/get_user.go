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
	GetUser func(ctx *fiber.Ctx) error
)

func NewGetUser(getUser query.GetUser) GetUser {
	return func(ctx *fiber.Ctx) error {
		userID, err := strconv.Atoi(ctx.Params("user_id"))
		if err != nil {
			return except.BadRequestError("invalid user id").SetInternal(err)
		}

		user, err := getUser(ctx.Context(), userID)
		if err != nil {
			return fmt.Errorf("query.getUser: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[*query.User](http.StatusOK, user))
	}
}
