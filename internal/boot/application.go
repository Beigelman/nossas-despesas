package boot

import (
	"context"

	"github.com/Beigelman/nossas-despesas/internal/controller/handler"
	"github.com/Beigelman/nossas-despesas/internal/controller/middleware"
	"github.com/Beigelman/nossas-despesas/internal/infra/postgres/authrepo"
	"github.com/Beigelman/nossas-despesas/internal/infra/postgres/categorygrouprepo"
	"github.com/Beigelman/nossas-despesas/internal/infra/postgres/categoryrepo"
	"github.com/Beigelman/nossas-despesas/internal/infra/postgres/expenserepo"
	"github.com/Beigelman/nossas-despesas/internal/infra/postgres/groupinviterepo"
	"github.com/Beigelman/nossas-despesas/internal/infra/postgres/grouprepo"
	"github.com/Beigelman/nossas-despesas/internal/infra/postgres/incomerepo"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/query"
	"github.com/Beigelman/nossas-despesas/internal/usecase"
)

var ApplicationModule = eon.NewModule("Application", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	// auth
	di.Provide(c, authrepo.NewPGRepository)
	di.Provide(c, usecase.NewSignUpWithCredentials)
	di.Provide(c, usecase.NewSignInWithCredentials)
	di.Provide(c, usecase.NewRefreshAuthToken)
	di.Provide(c, usecase.NewSignInWithGoogle)
	di.Provide(c, handler.NewSignUpWithCredentials)
	di.Provide(c, handler.NewSignInWithCredentials)
	di.Provide(c, handler.NewRefreshAuthToken)
	di.Provide(c, handler.NewSignInWithGoogle)
	di.Provide(c, middleware.NewAuthMiddleware)
	// income
	di.Provide(c, incomerepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateIncome)
	di.Provide(c, usecase.NewUpdateIncome)
	di.Provide(c, usecase.NewDeleteIncome)
	di.Provide(c, handler.NewCreateIncome)
	di.Provide(c, handler.NewUpdateIncome)
	di.Provide(c, handler.NewDeleteIncome)
	// group
	di.Provide(c, groupinviterepo.NewPGRepository)
	di.Provide(c, grouprepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateGroup)
	di.Provide(c, usecase.NewInviteUserToGroup)
	di.Provide(c, usecase.NewAcceptGroupInvite)
	di.Provide(c, query.NewGetGroup)
	di.Provide(c, query.NewGetGroupExpenses)
	di.Provide(c, query.NewGetGroupBalance)
	di.Provide(c, query.NewGetGroupMonthlyIncome)
	di.Provide(c, handler.NewInviteUserToGroup)
	di.Provide(c, handler.NewAcceptGroupInvite)
	di.Provide(c, handler.NewGetGroupBalance)
	di.Provide(c, handler.NewGetGroupExpenses)
	di.Provide(c, handler.NewGetGroupMonthlyIncome)
	di.Provide(c, handler.NewCreateGroup)
	di.Provide(c, handler.NewGetGroup)
	// expense
	di.Provide(c, expenserepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateExpense)
	di.Provide(c, usecase.NewUpdateExpense)
	di.Provide(c, usecase.NewDeleteExpense)
	di.Provide(c, usecase.NewRecalculateExpensesSplitRatio)
	di.Provide(c, query.NewGetExpensesPerSearch)
	di.Provide(c, query.NewGetExpenseDetails)
	di.Provide(c, handler.NewGetExpensesPerSearch)
	di.Provide(c, handler.NewCreateExpense)
	di.Provide(c, handler.NewUpdateExpense)
	di.Provide(c, handler.NewDeleteExpense)
	di.Provide(c, handler.NewGetExpenseDetails)
	// category
	di.Provide(c, categoryrepo.NewPGRepository)
	di.Provide(c, categorygrouprepo.NewPGRepository)
	di.Provide(c, usecase.NewCreateCategory)
	di.Provide(c, usecase.NewCreateCategoryGroup)
	di.Provide(c, query.NewGetCategories)
	di.Provide(c, handler.NewGetCategories)
	di.Provide(c, handler.NewCreateCategory)
	di.Provide(c, handler.NewCreateCategoryGroup)
	// Insights
	di.Provide(c, query.NewGetExpensesPerPeriod)
	di.Provide(c, query.NewGetExpensesPerCategory)
	di.Provide(c, query.NewGetIncomesPerPeriod)
	di.Provide(c, handler.NewGetExpensesPerPeriod)
	di.Provide(c, handler.NewGetExpensesPerCategory)
	di.Provide(c, handler.NewGetIncomesPerPeriod)
})
