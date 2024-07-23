package controller

import (
	"github.com/Beigelman/nossas-despesas/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(
	server *fiber.App,
	createExpenseHandler CreateExpense,
	getExpensesHandler GetExpenses,
	updateExpenseHandler UpdateExpense,
	deleteExpenseHandler DeleteExpense,
	authMiddleware middleware.AuthMiddleware,
	getExpensesPerCategoryHandler GetExpensesPerCategory,
	getExpensesPerPeriodHandler GetExpensesPerPeriod,
	getExpenseDetailsHandler GetExpenseDetails,
) {
	// Api group
	api := server.Group("api")
	// Api version V1
	v1 := api.Group("v1")
	// Expense routes
	expense := v1.Group("expenses", authMiddleware)
	expense.Post("/", createExpenseHandler)
	expense.Get("/", getExpensesHandler)
	expense.Get("/:expense_id/details", getExpenseDetailsHandler)
	expense.Patch("/:expense_id", updateExpenseHandler)
	expense.Delete("/:expense_id", deleteExpenseHandler)
	// Expenses insights routes
	insights := expense.Group("insights", authMiddleware)
	insights.Get("/", getExpensesPerPeriodHandler)
	insights.Get("/category", getExpensesPerCategoryHandler)
}
