package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/query"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type GetCategories func(ctx *fiber.Ctx) error

func NewGetCategories(getCategories query.GetCategories) GetCategories {
	return func(ctx *fiber.Ctx) error {
		categories, err := getCategories(ctx.Context())
		if err != nil {
			return fmt.Errorf("query.getCategories: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[[]query.CategoryGroup](http.StatusOK, categories))
	}
}
