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
	getExpensesPerCategoryHandler GetExpensesPerCategory,
	getExpensesPerPeriodHandler GetExpensesPerPeriod,
	getExpenseDetailsHandler GetExpenseDetails,
	generateExpensesFromScheduledHandler GenerateExpensesFromScheduled,
	createScheduledExpenseHandler CreateScheduledExpense,
	authMiddleware middleware.AuthMiddleware,
) {
	// Api group
	api := server.Group("api")
	// Api version V1
	v1 := api.Group("v1")
	// Expense routes
	expense := v1.Group("expenses")
	expense.Post("/", authMiddleware, createExpenseHandler)
	expense.Get("/", authMiddleware, getExpensesHandler)
	expense.Get("/:expense_id/details", authMiddleware, getExpenseDetailsHandler)
	expense.Patch("/:expense_id", authMiddleware, updateExpenseHandler)
	expense.Delete("/:expense_id", authMiddleware, deleteExpenseHandler)
	expense.Post("/scheduled", authMiddleware, createScheduledExpenseHandler)
	// Generate expenses from scheduled does not need auth
	expense.Post("/scheduled/generate", generateExpensesFromScheduledHandler)

	// Expenses insights routes
	insights := expense.Group("insights", authMiddleware)
	insights.Get("/", getExpensesPerPeriodHandler)
	insights.Get("/category", getExpensesPerCategoryHandler)
}
