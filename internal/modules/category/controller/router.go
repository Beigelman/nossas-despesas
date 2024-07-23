package controller

import (
	"github.com/Beigelman/nossas-despesas/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(
	server *fiber.App,
	createCategoryHandler CreateCategory,
	createCategoryGroupHandler CreateCategoryGroup,
	getCategoriesHandler GetCategories,
	authMiddleware middleware.AuthMiddleware,
) {
	// Api group
	api := server.Group("api")
	// Api version V1
	v1 := api.Group("v1")
	// Category routes
	category := v1.Group("category", authMiddleware)
	category.Get("/", getCategoriesHandler)
	category.Post("/", createCategoryHandler)
	category.Post("/group", createCategoryGroupHandler)
}
