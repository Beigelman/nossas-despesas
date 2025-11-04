package controller

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/modules/group/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type (
	AcceptGroupInvite func(ctx *fiber.Ctx) error

	AcceptGroupInviteResponse struct{}
)

func NewAcceptGroupInvite(acceptGroupInvite usecase.AcceptGroupInvite) AcceptGroupInvite {
	return func(ctx *fiber.Ctx) error {
		email, ok := ctx.Locals("email").(string)
		if !ok {
			return except.UnprocessableEntityError().SetInternal(fmt.Errorf("email not found in context"))
		}

		token := ctx.Params("token")

		if token == "" {
			return except.UnprocessableEntityError().SetInternal(fmt.Errorf("token not found in params"))
		}

		if err := acceptGroupInvite(ctx.Context(), usecase.AcceptGroupInviteInput{
			Email: email,
			Token: token,
		}); err != nil {
			return fmt.Errorf("usecase.AcceptGroupInvite: %w", err)
		}

		return ctx.Status(http.StatusOK).SendString("Invite accepted successfully!")
	}
}
