package controller

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/group/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
)

type (
	InviteUserToGroup func(ctx *fiber.Ctx) error

	InviteUserToGroupRequest struct {
		Email   string `json:"email" validate:"required"`
		BaseURL string `json:"base_url" validate:"required"`
	}

	InviteUserToGroupResponse struct {
		Link string `json:"link"`
	}
)

func NewInviteUserToGroup(inviteUserToGroup usecase.InviteUserToGroup) InviteUserToGroup {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var request InviteUserToGroupRequest
		if err := ctx.BodyParser(&request); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(request); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.UnprocessableEntityError().SetInternal(fmt.Errorf("group_id not found in context"))
		}

		invite, err := inviteUserToGroup(ctx.Context(), usecase.InviteUserToGroupInput{
			GroupID: group.ID{Value: groupID},
			Email:   request.Email,
			BaseURL: request.BaseURL,
		})
		if err != nil {
			return fmt.Errorf("usecase.InviteUserToGroup: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, invite))
	}
}
