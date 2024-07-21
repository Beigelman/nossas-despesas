package controller

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/group/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"net/http"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/gofiber/fiber/v2"
)

type (
	CreateGroup func(ctx *fiber.Ctx) error

	CreateGroupRequest struct {
		Name string `json:"name" validate:"required,min=3,max=50"`
	}

	CreateGroupResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

func NewCreateGroup(createGroup usecase.CreateGroup) CreateGroup {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req CreateGroupRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		userID, ok := ctx.Locals("user_id").(int)
		if !ok {
			return except.UnprocessableEntityError().SetInternal(fmt.Errorf("user_id not found in context"))
		}

		group, err := createGroup(ctx.Context(), usecase.CreateGroupInput{
			Name:   req.Name,
			UserID: user.ID{Value: userID},
		})
		if err != nil {
			return fmt.Errorf("createGroup: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, CreateGroupResponse{
				ID:   group.ID.Value,
				Name: group.Name,
			}),
		)
	}
}
