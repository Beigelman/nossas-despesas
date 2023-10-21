package boot

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/categorygrouprepo"

	"github.com/Beigelman/ludaapi/internal/controller/handler"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/categoryrepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/expenserepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/grouprepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/userrepo"
	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/Beigelman/ludaapi/internal/pkg/eon"
	"github.com/Beigelman/ludaapi/internal/usecase"
)

var ApplicationModule = eon.NewModule("Application", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	// user
	di.Provide(c, userrepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateUser)
	di.Provide(c, handler.NewCreateUserHandler)
	// group
	di.Provide(c, grouprepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateGroup)
	di.Provide(c, handler.NewCreateGroupHandler)
	// expense
	di.Provide(c, expenserepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateExpense)
	di.Provide(c, handler.NewCreateExpenseHandler)
	// category
	di.Provide(c, categoryrepo.NewPGRepository)
	di.Provide(c, categorygrouprepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateCategory)
	di.Provide(c, usecase.NewCreateCategoryGroup)
	di.Provide(c, handler.NewCreateCategoryHandler)
	di.Provide(c, handler.NewCreateCategoryGroupHandler)
})
