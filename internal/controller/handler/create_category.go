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

func NewCreateCategoryHandler(createCategory usecase.CreateCategory) CreateCategory {
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
			CategoryGroupID: entity.CategoryGroupID{Value: req.CategoryGroupID},
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
