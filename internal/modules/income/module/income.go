package income

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/postgres"

	"github.com/Beigelman/nossas-despesas/internal/modules/income/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/query"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var Module = eon.NewModule("Income", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	// income
	di.Provide(c, postgres.NewIncomeRepository)
	di.Provide(c, usecase.NewCreateIncome)
	di.Provide(c, usecase.NewUpdateIncome)
	di.Provide(c, usecase.NewDeleteIncome)
	di.Provide(c, query.NewGetIncomesPerPeriod)
	di.Provide(c, query.NewGetMonthlyIncome)
	di.Provide(c, controller.NewCreateIncome)
	di.Provide(c, controller.NewUpdateIncome)
	di.Provide(c, controller.NewDeleteIncome)
	di.Provide(c, controller.NewGetMonthlyIncome)
	di.Provide(c, controller.NewGetIncomesPerPeriod)
	// Register routes
	lc.OnBooted(eon.HookOrders.APPEND, func() error { return di.Call(c, controller.Router) })
})
