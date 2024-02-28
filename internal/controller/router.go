package controller

import (
	"github.com/Beigelman/nossas-despesas/internal/controller/handler"
	"github.com/Beigelman/nossas-despesas/internal/controller/middleware"
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
	getMyUserHandler handler.GetMe,
	getCategoriesHandler handler.GetCategories,
	updateExpenseHandler handler.UpdateExpense,
	deleteExpenseHandler handler.DeleteExpense,
	getGroupBalanceHandler handler.GetGroupBalance,
	createIncomeHandler handler.CreateIncome,
	updateIncomeHandler handler.UpdateIncome,
	deleteIncomeHandler handler.DeleteIncome,
	getGroupMonthlyIncomeHandler handler.GetGroupMonthlyIncome,
	signInWithCredentialsHandler handler.SignInWithCredentials,
	signInWithGoogleHandler handler.SignInWithGoogle,
	signUpWithCredentialsHandler handler.SignUpWithCredentials,
	refreshAuthTokenHandler handler.RefreshAuthToken,
	authMiddleware middleware.AuthMiddleware,
	getExpensesPerCategoryHandler handler.GetExpensesPerCategory,
	getExpensesPerPeriodHandler handler.GetExpensesPerPeriod,
	getIncomesPerPeriodHandler handler.GetIncomesPerPeriod,
	getExpenseDetailsHandler handler.GetExpenseDetails,
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
		auth.Post("/sign-in/google", signInWithGoogleHandler)
		auth.Post("/sign-up/credentials", signUpWithCredentialsHandler)
		auth.Post("refresh-token", refreshAuthTokenHandler)
		// User routes
		user := v1.Group("user", authMiddleware)
		user.Get("/me", getMyUserHandler)
		user.Patch("/add-to-group", addUserToGroupHandler)
		// Group routes
		group := v1.Group("group", authMiddleware)
		group.Post("/", createGroupHandler)
		group.Get("/", getGroupHandler)
		group.Get("/expenses", getGroupExpenseHandler)
		group.Get("/balance", getGroupBalanceHandler)
		group.Get("income", getGroupMonthlyIncomeHandler)
		// Expense routes
		expense := v1.Group("expense", authMiddleware)
		expense.Post("/", createExpenseHandler)
		expense.Get("/:expense_id/details", getExpenseDetailsHandler)
		expense.Patch("/:expense_id", updateExpenseHandler)
		expense.Delete("/:expense_id", deleteExpenseHandler)
		// Income routes
		income := v1.Group("income", authMiddleware)
		income.Post("/", createIncomeHandler)
		income.Patch("/:income_id", updateIncomeHandler)
		income.Delete("/:income_id", deleteIncomeHandler)
		// Category routes
		category := v1.Group("category", authMiddleware)
		category.Get("/", getCategoriesHandler)
		category.Post("/", createCategoryHandler)
		category.Post("/group", createCategoryGroupHandler)
		// Insights routes
		insights := v1.Group("insights", authMiddleware)
		insights.Get("/expenses", getExpensesPerPeriodHandler)
		insights.Get("/expenses/category", getExpensesPerCategoryHandler)
		insights.Get("/incomes", getIncomesPerPeriodHandler)

	}
}
