package controller

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/category/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
)

type (
	CreateCategory func(ctx *fiber.Ctx) error

	CreateCategoryRequest struct {
		Icon            string `json:"icon" validate:"required"`
		Name            string `json:"name" validate:"required"`
		CategoryGroupID int    `json:"category_group_id" validate:"required"`
	}

	CreateCategoryResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

func NewCreateCategory(createCategory usecase.CreateCategory) CreateCategory {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req CreateCategoryRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		group, err := createCategory(ctx.Context(), usecase.CreateCategoryInput{
			Name:            req.Name,
			Icon:            req.Icon,
			CategoryGroupID: category.GroupID{Value: req.CategoryGroupID},
		})
		if err != nil {
			return fmt.Errorf("createCategory: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, CreateCategoryResponse{
				ID:   group.ID.Value,
				Name: group.Name,
			}),
		)
	}
}
