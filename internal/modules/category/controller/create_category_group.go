package controller

import (
	"fmt"
	"net/http"

	"github.com/Beigelman/nossas-despesas/internal/modules/category/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type (
	CreateCategoryGroup func(ctx *fiber.Ctx) error

	CreateCategoryGroupRequest struct {
		Icon string `json:"icon" validate:"required"`
		Name string `json:"name" validate:"required"`
	}

	CreateCategoryGroupResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

func NewCreateCategoryGroup(createCategoryGroup usecase.CreateCategoryGroup) CreateCategoryGroup {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req CreateCategoryGroupRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		categoryGroup, err := createCategoryGroup(ctx.Context(), usecase.CreateCategoryGroupInput{
			Name: req.Name,
			Icon: req.Icon,
		})
		if err != nil {
			return fmt.Errorf("createCategoryGroup: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, CreateCategoryGroupResponse{
				ID:   categoryGroup.ID.Value,
				Name: categoryGroup.Name,
			}),
		)
	}
}
