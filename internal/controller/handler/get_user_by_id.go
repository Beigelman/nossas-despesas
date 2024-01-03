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
	GetUserByID func(ctx *fiber.Ctx) error
)

func NewGetUserByID(getUserID query.GetUserByID) GetUserByID {
	return func(ctx *fiber.Ctx) error {
		userID, err := strconv.Atoi(ctx.Params("user_id"))
		if err != nil {
			return except.BadRequestError("invalid user_id").SetInternal(err)
		}

		user, err := getUserID(ctx.Context(), userID)
		if err != nil {
			return fmt.Errorf("query.getUser: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[*query.User](http.StatusOK, user))
	}
}
