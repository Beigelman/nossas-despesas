package handler

import (
	"fmt"
	"net/http"

	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type (
	CreateGroup        func(ctx *fiber.Ctx) error
	CreateGroupRequest struct {
		Name string `json:"name"`
	}
	CreateGroupResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

func NewCreateGroupHandler(createGroup usecase.CreateGroup) CreateGroup {
	return func(ctx *fiber.Ctx) error {
		var req CreateGroupRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("ctx.BodyParser: %w", err)
		}

		group, err := createGroup(ctx.Context(), req.Name)
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
