package category

import (
	"context"

	"github.com/Beigelman/nossas-despesas/internal/modules/category/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/category/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/category/query"
	"github.com/Beigelman/nossas-despesas/internal/modules/category/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var Module = eon.NewModule("Category", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, postgres.NewCategoryRepository)
	di.Provide(c, postgres.NewCategoryGroupRepository)
	di.Provide(c, usecase.NewCreateCategory)
	di.Provide(c, usecase.NewCreateCategoryGroup)
	di.Provide(c, query.NewGetCategories)
	di.Provide(c, controller.NewGetCategories)
	di.Provide(c, controller.NewCreateCategory)
	di.Provide(c, controller.NewCreateCategoryGroup)
	// Register routes
	lc.OnBooted(eon.HookOrders.APPEND, func() error { return di.Call(c, controller.Router) })
})
