package boot

import (
	"context"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/infra/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/query"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var ExpenseModule = eon.NewModule("Expense", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	// expense
	di.Provide(c, postgres.NewExpenseRepository)
	di.Provide(c, usecase.NewCreateExpense)
	di.Provide(c, usecase.NewUpdateExpense)
	di.Provide(c, usecase.NewDeleteExpense)
	di.Provide(c, usecase.NewRecalculateExpensesSplitRatio)
	di.Provide(c, query.NewGetExpenses)
	di.Provide(c, query.NewGetExpenseDetails)
	di.Provide(c, query.NewGetExpensesPerPeriod)
	di.Provide(c, query.NewGetExpensesPerCategory)
	di.Provide(c, controller.NewGetExpenses)
	di.Provide(c, controller.NewCreateExpense)
	di.Provide(c, controller.NewUpdateExpense)
	di.Provide(c, controller.NewDeleteExpense)
	di.Provide(c, controller.NewGetExpenseDetails)
	di.Provide(c, controller.NewGetExpensesPerPeriod)
	di.Provide(c, controller.NewGetExpensesPerCategory)
	di.Provide(c, controller.NewRecalculateExpensesSplitRatio)
	// Register routes
	lc.OnBooted(eon.HookOrders.APPEND, func() error {
		return di.Call(c, controller.Router)
	})
	// Listen to subscriber
	lc.OnRunning(eon.HookOrders.APPEND, func() error {
		recalculate := di.Resolve[controller.RecalculateExpensesSplitRatio](c)
		return recalculate(ctx)
	})
})
