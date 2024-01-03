package controller

import (
	"firebase.google.com/go/v4/auth"
	"github.com/Beigelman/ludaapi/internal/controller/handler"
	"github.com/Beigelman/ludaapi/internal/controller/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(
	server *fiber.App,
	auth *auth.Client,
	createUserHandler handler.CreateUser,
	createGroupHandler handler.CreateGroup,
	createExpenseHandler handler.CreateExpense,
	createCategoryHandler handler.CreateCategory,
	createCategoryGroupHandler handler.CreateCategoryGroup,
	addUserToGroupHandler handler.AddUserToGroup,
	getGroupExpenseHandler handler.GetGroupExpenses,
	getGroupHandler handler.GetGroup,
	getUserByAuthenticationIDHandler handler.GetUserByAuthenticationID,
	getUserByIDHandler handler.GetUserByID,
	getCategoriesHandler handler.GetCategories,
	updateExpenseHandler handler.UpdateExpense,
	deleteExpenseHandler handler.DeleteExpense,
	getGroupBalanceHandler handler.GetGroupBalance,
) {
	server.Get("healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	authMiddleware := middleware.AuthMiddleware(auth)
	// Api group
	api := server.Group("api")
	// Api version V1
	{
		v1 := api.Group("v1")
		// User routes
		user := v1.Group("user")
		user.Post("/", createUserHandler)
		user.Get("/:user_id", getUserByIDHandler, authMiddleware)
		user.Get("/authentication/:authentication_id", getUserByAuthenticationIDHandler)
		user.Patch("/add-to-group", addUserToGroupHandler, authMiddleware)
		// Group routes
		group := v1.Group("group", authMiddleware)
		group.Post("/", createGroupHandler)
		group.Get("/:group_id", getGroupHandler)
		group.Get("/:group_id/expenses", getGroupExpenseHandler)
		group.Get("/:group_id/balance", getGroupBalanceHandler)
		// Expense routes
		expense := v1.Group("expense", authMiddleware)
		expense.Post("/", createExpenseHandler)
		expense.Patch("/:expense_id", updateExpenseHandler)
		expense.Delete("/:expense_id", deleteExpenseHandler)
		// Category routes
		category := v1.Group("category", authMiddleware)
		category.Get("/", getCategoriesHandler)
		category.Post("/", createCategoryHandler)
		category.Post("/group", createCategoryGroupHandler)
	}
}
