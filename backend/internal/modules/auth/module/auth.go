package auth

import (
	"context"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var Module = eon.NewModule("Auth", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, postgres.NewAuthRepository)
	di.Provide(c, usecase.NewSignUpWithCredentials)
	di.Provide(c, usecase.NewSignInWithCredentials)
	di.Provide(c, usecase.NewRefreshAuthToken)
	di.Provide(c, usecase.NewSignInWithGoogle)
	di.Provide(c, controller.NewSignUpWithCredentials)
	di.Provide(c, controller.NewSignInWithCredentials)
	di.Provide(c, controller.NewRefreshAuthToken)
	di.Provide(c, controller.NewSignInWithGoogle)
	// Register Routes
	lc.OnBooted(eon.HookOrders.APPEND, func() error { return di.Call(c, controller.Router) })
})
