package boot

import (
	"context"

	"github.com/Beigelman/nossas-despesas/internal/modules/user/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/user/controller/handler"
	"github.com/Beigelman/nossas-despesas/internal/modules/user/infra/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/user/query"
	"github.com/Beigelman/nossas-despesas/internal/modules/user/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var UserModule = eon.NewModule("User", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, postgres.NewUserRepository)
	di.Provide(c, usecase.NewCreateUser)
	di.Provide(c, query.NewGetUserByID)
	di.Provide(c, handler.NewGetMe)

	lc.OnBooted(eon.HookOrders.PREPEND, func() error {
		return di.Call(c, controller.Router)
	})
})
