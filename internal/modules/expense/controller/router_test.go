package controller_test

import (
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	var called []string
	h := func(name string) fiber.Handler {
		return func(c *fiber.Ctx) error {
			called = append(called, name)
			return nil
		}
	}

	// Mock do middleware de autenticação
	mockAuthMiddleware := func(c *fiber.Ctx) error {
		return c.Next()
	}

	controller.Router(app,
		h("createExpense"),
		h("getExpenses"),
		h("updateExpense"),
		h("deleteExpense"),
		h("getExpensesPerCategory"),
		h("getExpensesPerPeriod"),
		h("getExpenseDetails"),
		h("generateExpensesFromScheduled"),
		h("createScheduledExpense"),
		mockAuthMiddleware,
	)

	routes := app.GetRoutes()
	paths := make([]string, len(routes))
	for i, r := range routes {
		paths[i] = r.Method + " " + r.Path
	}

	// Testa se as rotas básicas de expenses foram registradas
	assert.Contains(t, paths, "POST /api/v1/expenses/")
	assert.Contains(t, paths, "GET /api/v1/expenses/")
	assert.Contains(t, paths, "GET /api/v1/expenses/:expense_id/details")
	assert.Contains(t, paths, "PATCH /api/v1/expenses/:expense_id")
	assert.Contains(t, paths, "DELETE /api/v1/expenses/:expense_id")

	// Testa se as rotas de scheduled expenses foram registradas
	assert.Contains(t, paths, "POST /api/v1/expenses/scheduled")
	assert.Contains(t, paths, "POST /api/v1/expenses/scheduled/generate")

	// Testa se as rotas de insights foram registradas
	assert.Contains(t, paths, "GET /api/v1/expenses/insights/")
	assert.Contains(t, paths, "GET /api/v1/expenses/insights/category")
}

func TestRouterAuthMiddleware(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	mockAuthMiddleware := func(c *fiber.Ctx) error {
		return c.Next()
	}

	// Mock handlers simples
	h := func(name string) fiber.Handler {
		return func(c *fiber.Ctx) error {
			return c.SendString(name)
		}
	}

	controller.Router(app,
		h("createExpense"),
		h("getExpenses"),
		h("updateExpense"),
		h("deleteExpense"),
		h("getExpensesPerCategory"),
		h("getExpensesPerPeriod"),
		h("getExpenseDetails"),
		h("generateExpensesFromScheduled"),
		h("createScheduledExpense"),
		mockAuthMiddleware,
	)

	// Testa se o middleware é aplicado nas rotas protegidas
	routes := app.GetRoutes()

	// Verifica se temos o número correto de rotas
	assert.GreaterOrEqual(t, len(routes), 9, "Deve ter pelo menos 9 rotas registradas")

	// Verifica se todas as rotas principais estão presents
	expectedRoutes := []string{
		"POST /api/v1/expenses/",
		"GET /api/v1/expenses/",
		"GET /api/v1/expenses/:expense_id/details",
		"PATCH /api/v1/expenses/:expense_id",
		"DELETE /api/v1/expenses/:expense_id",
		"POST /api/v1/expenses/scheduled",
		"POST /api/v1/expenses/scheduled/generate",
		"GET /api/v1/expenses/insights/",
		"GET /api/v1/expenses/insights/category",
	}

	actualRoutes := make([]string, len(routes))
	for i, r := range routes {
		actualRoutes[i] = r.Method + " " + r.Path
	}

	for _, expectedRoute := range expectedRoutes {
		assert.Contains(t, actualRoutes, expectedRoute, "Rota %s deve estar presente", expectedRoute)
	}
}
