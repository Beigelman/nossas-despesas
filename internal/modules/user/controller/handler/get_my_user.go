package handler

import (
	"fmt"
	"net/http"

	"github.com/Beigelman/nossas-despesas/internal/modules/user/query"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
)

type (
	GetMe func(ctx *fiber.Ctx) error
)

func NewGetMe(getUserID query.GetUserByID) GetMe {
	return func(ctx *fiber.Ctx) error {
		userID, ok := ctx.Locals("user_id").(int)
		if !ok {
			return except.BadRequestError("invalid user_id")
		}

		user, err := getUserID(ctx.Context(), userID)
		if err != nil {
			return fmt.Errorf("query.getUser: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, user))
	}
}
