package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/query"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type (
	GetUserByAuthenticationID func(ctx *fiber.Ctx) error
)

func NewGetUserByAuthenticationID(getUserByAuthenticationID query.GetUserByAuthenticationID) GetUserByAuthenticationID {
	return func(ctx *fiber.Ctx) error {
		authenticationID := ctx.Params("authentication_id")
		if authenticationID == "" {
			return fmt.Errorf("authentication_id is required")
		}

		user, err := getUserByAuthenticationID(ctx.Context(), authenticationID)
		if err != nil {
			return fmt.Errorf("query.getUser: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[*query.User](http.StatusOK, user))
	}
}
