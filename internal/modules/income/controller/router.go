package controller

import (
	"github.com/Beigelman/nossas-despesas/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(
	server *fiber.App,
	createIncomeHandler CreateIncome,
	updateIncomeHandler UpdateIncome,
	deleteIncomeHandler DeleteIncome,
	authMiddleware middleware.AuthMiddleware,
	getIncomesPerPeriodHandler GetIncomesPerPeriod,

) {

	// Api group
	api := server.Group("api")
	// Api version V1
	{
		v1 := api.Group("v1")
		// Income routes
		income := v1.Group("income", authMiddleware)
		income.Post("/", createIncomeHandler)
		income.Patch("/:income_id", updateIncomeHandler)
		income.Delete("/:income_id", deleteIncomeHandler)
		income.Get("/insights", getIncomesPerPeriodHandler)
	}
}
