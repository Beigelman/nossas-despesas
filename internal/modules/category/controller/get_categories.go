package controller

import (
	"fmt"
	"net/http"

	"github.com/Beigelman/nossas-despesas/internal/modules/category/query"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/gofiber/fiber/v2"
)

type GetCategories func(ctx *fiber.Ctx) error

func NewGetCategories(getCategories query.GetCategories) GetCategories {
	return func(ctx *fiber.Ctx) error {
		categories, err := getCategories(ctx.Context())
		if err != nil {
			return fmt.Errorf("query.getCategories: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, categories))
	}
}
