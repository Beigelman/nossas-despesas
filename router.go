package main

//import (
//	"github.com/Beigelman/nossas-despesas/internal/modules/auth/controller"
//	controller2 "github.com/Beigelman/nossas-despesas/internal/modules/category/controller"
//	controller5 "github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
//	controller3 "github.com/Beigelman/nossas-despesas/internal/modules/group/controller"
//	controller4 "github.com/Beigelman/nossas-despesas/internal/modules/income/controller"
//	"github.com/Beigelman/nossas-despesas/internal/pkg/middleware"
//	"github.com/gofiber/fiber/v2"
//)
//
//func Router(
//	server *fiber.App,
//	createGroupHandler controller3.CreateGroup,
//	createExpenseHandler controller5.CreateExpense,
//	createCategoryHandler controller2.CreateCategory,
//	createCategoryGroupHandler controller2.CreateCategoryGroup,
//	getGroupExpenseHandler controller5.GetGroupExpenses,
//	getExpensesPerSearch controller5.GetExpensesPerSearch,
//	getGroupHandler controller3.GetGroup,
//	getCategoriesHandler controller2.GetCategories,
//	updateExpenseHandler controller5.UpdateExpense,
//	deleteExpenseHandler controller5.DeleteExpense,
//	getGroupBalanceHandler controller3.GetGroupBalance,
//	createIncomeHandler controller4.CreateIncome,
//	updateIncomeHandler controller4.UpdateIncome,
//	deleteIncomeHandler controller4.DeleteIncome,
//	getGroupMonthlyIncomeHandler controller4.GetGroupMonthlyIncome,
//	signInWithCredentialsHandler controller.SignInWithCredentials,
//	signUpWithCredentialsHandler controller.SignUpWithCredentials,
//	signInWithGoogleHandler controller.SignInWithGoogle,
//	refreshAuthTokenHandler controller.RefreshAuthToken,
//	authMiddleware middleware.AuthMiddleware,
//	getExpensesPerCategoryHandler controller5.GetExpensesPerCategory,
//	getExpensesPerPeriodHandler controller5.GetExpensesPerPeriod,
//	getIncomesPerPeriodHandler controller4.GetIncomesPerPeriod,
//	getExpenseDetailsHandler controller5.GetExpenseDetails,
//	inviteUserToGroupHandler controller3.InviteUserToGroup,
//	acceptGroupInviteHandler controller3.AcceptGroupInvite,
//) {
//
//	// Api group
//	api := server.Group("api")
//	// Api version V1
//	{
//		v1 := api.Group("v1")
//		// Auth routes
//		auth := v1.Group("auth")
//		auth.Post("/sign-in/credentials", signInWithCredentialsHandler)
//		auth.Post("/sign-in/google", signInWithGoogleHandler)
//		auth.Post("/sign-up/credentials", signUpWithCredentialsHandler)
//		auth.Post("refresh-token", refreshAuthTokenHandler)
//		// Group routes
//		group := v1.Group("group", authMiddleware)
//		group.Post("/", createGroupHandler)
//		group.Post("/invite", inviteUserToGroupHandler)
//		group.Post("/invite/:token/accept", acceptGroupInviteHandler)
//		group.Get("/", getGroupHandler)
//		group.Get("/balance", getGroupBalanceHandler)
//		group.Get("/expenses", getGroupExpenseHandler)
//		group.Get("income", getGroupMonthlyIncomeHandler)
//		// Expense routes
//		expense := v1.Group("expense", authMiddleware)
//		expense.Post("/", createExpenseHandler)
//		expense.Get("/", getExpensesPerSearch)
//		expense.Get("/:expense_id/details", getExpenseDetailsHandler)
//		expense.Patch("/:expense_id", updateExpenseHandler)
//		expense.Delete("/:expense_id", deleteExpenseHandler)
//		// Income routes
//		income := v1.Group("income", authMiddleware)
//		income.Post("/", createIncomeHandler)
//		income.Patch("/:income_id", updateIncomeHandler)
//		income.Delete("/:income_id", deleteIncomeHandler)
//		// Category routes
//		category := v1.Group("category", authMiddleware)
//		category.Get("/", getCategoriesHandler)
//		category.Post("/", createCategoryHandler)
//		category.Post("/group", createCategoryGroupHandler)
//		// Insights routes
//		insights := v1.Group("insights", authMiddleware)
//		insights.Get("/expenses", getExpensesPerPeriodHandler)
//		insights.Get("/expenses/category", getExpensesPerCategoryHandler)
//		insights.Get("/incomes", getIncomesPerPeriodHandler)
//	}
//}
