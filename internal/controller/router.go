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
		user.Patch("/add-to-group", addUserToGroupHandler)
		// Group routes
		group := v1.Group("group")
		group.Post("/", createGroupHandler)
		// Expense routes
		expense := v1.Group("expense")
		expense.Post("/", createExpenseHandler)
		// Expenses routes
		expenses := v1.Group("expenses")
		expenses.Get("/group/:group_id", getGroupExpenseHandler)
		// Category routes
		category := v1.Group("category")
		category.Post("/", createCategoryHandler)
		category.Post("/group", createCategoryGroupHandler)
	}
}
