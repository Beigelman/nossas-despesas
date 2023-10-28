package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/pkg/validator"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type (
	AddUserToGroup func(ctx *fiber.Ctx) error

	AddUserToGroupRequest struct {
		GroupID int `json:"group_id" validate:"required"`
		UserID  int `json:"user_id" validate:"required"`
	}

	AddUserToGroupResponse struct{}
)

func NewAddUserToGroup(addUserToGroup usecase.AddUserToGroup) AddUserToGroup {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var request AddUserToGroupRequest
		if err := ctx.BodyParser(&request); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(request); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		user, err := addUserToGroup(ctx.Context(), usecase.AddUserToGroupInput{
			GroupID: entity.GroupID{Value: request.GroupID},
			UserID:  entity.UserID{Value: request.UserID},
		})
		if err != nil {
			return fmt.Errorf("usecase.AddUserToGroup: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[*entity.User](http.StatusOK, user))
	}
}
