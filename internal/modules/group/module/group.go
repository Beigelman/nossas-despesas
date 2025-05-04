package group

import (
	"context"

	"github.com/Beigelman/nossas-despesas/internal/modules/group/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/group/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/group/query"
	"github.com/Beigelman/nossas-despesas/internal/modules/group/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var Module = eon.NewModule("Group", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	// group
	di.Provide(c, postgres.NewGroupRepository)
	di.Provide(c, postgres.NewGroupInviteRepository)
	di.Provide(c, usecase.NewCreateGroup)
	di.Provide(c, usecase.NewInviteUserToGroup)
	di.Provide(c, usecase.NewAcceptGroupInvite)
	di.Provide(c, query.NewGetGroup)
	di.Provide(c, query.NewGetGroupBalance)
	di.Provide(c, controller.NewInviteUserToGroup)
	di.Provide(c, controller.NewAcceptGroupInvite)
	di.Provide(c, controller.NewGetGroupBalance)
	di.Provide(c, controller.NewCreateGroup)
	di.Provide(c, controller.NewGetGroup)
	// Register routes
	lc.OnBooted(eon.HookOrders.APPEND, func() error { return di.Call(c, controller.Router) })
})
