package boot

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/categorygrouprepo"
	"github.com/Beigelman/ludaapi/internal/query"

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
	di.Provide(c, handler.NewCreateUser)
	di.Provide(c, query.NewGetUserByID)
	di.Provide(c, query.NewGetUserByAuthenticationID)
	di.Provide(c, usecase.NewAddUserToGroup)
	di.Provide(c, handler.NewAddUserToGroup)
	di.Provide(c, handler.NewGetUserByAuthenticationID)
	di.Provide(c, handler.NewGetUserByID)
	// group
	di.Provide(c, grouprepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateGroup)
	di.Provide(c, query.NewGetGroup)
	di.Provide(c, query.NewGetGroupExpenses)
	di.Provide(c, query.NewGetGroupBalance)
	di.Provide(c, handler.NewGetGroupBalance)
	di.Provide(c, handler.NewGetGroupExpenses)
	di.Provide(c, handler.NewCreateGroup)
	di.Provide(c, handler.NewGetGroup)
	// expense
	di.Provide(c, expenserepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateExpense)
	di.Provide(c, usecase.NewUpdateExpense)
	di.Provide(c, usecase.NewDeleteExpense)
	di.Provide(c, handler.NewCreateExpense)
	di.Provide(c, handler.NewUpdateExpense)
	di.Provide(c, handler.NewDeleteExpense)
	// category
	di.Provide(c, categoryrepo.NewPGRepository)
	di.Provide(c, categorygrouprepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateCategory)
	di.Provide(c, usecase.NewCreateCategoryGroup)
	di.Provide(c, query.NewGetCategories)
	di.Provide(c, handler.NewGetCategories)
	di.Provide(c, handler.NewCreateCategory)
	di.Provide(c, handler.NewCreateCategoryGroup)
})
