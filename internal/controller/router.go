package controller

import (
	"github.com/Beigelman/ludaapi/internal/controller/handler"
	"github.com/gofiber/fiber/v2"
)

func Router(
	server *fiber.App,
	createUserHandler handler.CreateUser,
	createGroupHandler handler.CreateGroup,
	createExpenseHandler handler.CreateExpense,
	createCategoryHandler handler.CreateCategory,
	createCategoryGroupHandler handler.CreateCategoryGroup,
	addUserToGroupHandler handler.AddUserToGroup,
	getGroupExpenseHandler handler.GetGroupExpenses,
	getGroupHandler handler.GetGroup,
	getUserHandler handler.GetUser,
	getCategoriesHandler handler.GetCategories,
) {
	server.Get("healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	// Api group
	api := server.Group("api")
	// Api version V1
	{
		v1 := api.Group("v1")
		// User routes
		user := v1.Group("user")
		user.Post("/", createUserHandler)
		user.Get("/:user_id", getUserHandler)
		user.Patch("/add-to-group", addUserToGroupHandler)
		// Group routes
		group := v1.Group("group")
		group.Post("/", createGroupHandler)
		group.Get("/:group_id", getGroupHandler)
		group.Get("/:group_id/expenses", getGroupExpenseHandler)
		// Expense routes
		expense := v1.Group("expense")
		expense.Post("/", createExpenseHandler)
		// Category routes
		category := v1.Group("category")
		category.Get("/", getCategoriesHandler)
		category.Post("/", createCategoryHandler)
		category.Post("/group", createCategoryGroupHandler)
	}
}
