package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
)

type GenerateExpensesFromScheduled func(ctx *fiber.Ctx) error

func NewGenerateExpensesFromScheduled(generateExpensesFromScheduled usecase.GenerateExpensesFromScheduledUseCase) GenerateExpensesFromScheduled {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		expensesCreated, err := generateExpensesFromScheduled(ctx)
		if err != nil {
			return fmt.Errorf("GenerateExpensesFromScheduled: %w", err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"expenses_created": expensesCreated,
		})
	}
}
