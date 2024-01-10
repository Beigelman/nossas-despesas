package controller

import (
	"github.com/Beigelman/ludaapi/internal/controller/handler"
	"github.com/Beigelman/ludaapi/internal/controller/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(
	server *fiber.App,
	createGroupHandler handler.CreateGroup,
	createExpenseHandler handler.CreateExpense,
	createCategoryHandler handler.CreateCategory,
	createCategoryGroupHandler handler.CreateCategoryGroup,
	addUserToGroupHandler handler.AddUserToGroup,
	getGroupExpenseHandler handler.GetGroupExpenses,
	getGroupHandler handler.GetGroup,
	getMyUserHandler handler.GetMyUser,
	getCategoriesHandler handler.GetCategories,
	updateExpenseHandler handler.UpdateExpense,
	deleteExpenseHandler handler.DeleteExpense,
	getGroupBalanceHandler handler.GetGroupBalance,
	signInWithCredentialsHandler handler.SignInWithCredentials,
	signUpWithCredentialsHandler handler.SignUpWithCredentials,
	refreshAuthTokenHandler handler.RefreshAuthToken,
	authMiddleware middleware.AuthMiddleware,
) {
	server.Get("healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Api group
	api := server.Group("api")
	// Api version V1
	{
		v1 := api.Group("v1")
		// Auth routes
		auth := v1.Group("auth")
		auth.Post("/sign-in/credentials", signInWithCredentialsHandler)
		auth.Post("/sign-up/credentials", signUpWithCredentialsHandler)
		auth.Post("refresh-jwt", refreshAuthTokenHandler)
		// User routes
		user := v1.Group("user", authMiddleware)
		user.Get("/me", getMyUserHandler)
		user.Patch("/add-to-group", addUserToGroupHandler)
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
