package main

//import (
//	"context"
//	"github.com/Beigelman/nossas-despesas/internal/modules/auth/controller"
//	postgres3 "github.com/Beigelman/nossas-despesas/internal/modules/auth/infra/postgres"
//	usecase4 "github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
//	controller5 "github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
//	postgres5 "github.com/Beigelman/nossas-despesas/internal/modules/expense/infra/postgres"
//	query5 "github.com/Beigelman/nossas-despesas/internal/modules/expense/query"
//	usecase6 "github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
//	controller3 "github.com/Beigelman/nossas-despesas/internal/modules/group/controller"
//	"github.com/Beigelman/nossas-despesas/internal/modules/group/infra/postgres"
//	query2 "github.com/Beigelman/nossas-despesas/internal/modules/group/query"
//	usecase2 "github.com/Beigelman/nossas-despesas/internal/modules/group/usecase"
//	controller4 "github.com/Beigelman/nossas-despesas/internal/modules/income/controller"
//	query4 "github.com/Beigelman/nossas-despesas/internal/modules/income/query"
//	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
//	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
//)
//
//var ApplicationModule = eon.NewModule("Application", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
//	// auth
//	di.Provide(c, postgres3.NewAuthRepository)
//	di.Provide(c, usecase4.NewSignUpWithCredentials)
//	di.Provide(c, usecase4.NewSignInWithCredentials)
//	di.Provide(c, usecase4.NewRefreshAuthToken)
//	di.Provide(c, usecase4.NewSignInWithGoogle)
//	di.Provide(c, controller.NewSignUpWithCredentials)
//	di.Provide(c, controller.NewSignInWithCredentials)
//	di.Provide(c, controller.NewRefreshAuthToken)
//	di.Provide(c, controller.NewSignInWithGoogle)
//
//	// group
//	di.Provide(c, postgres.NewGroupRepository)
//	di.Provide(c, postgres.NewGroupRepository)
//	di.Provide(c, usecase2.NewCreateGroup)
//	di.Provide(c, usecase2.NewInviteUserToGroup)
//	di.Provide(c, usecase2.NewAcceptGroupInvite)
//	di.Provide(c, query2.NewGetGroup)
//	di.Provide(c, query5.NewGetGroupExpenses)
//	di.Provide(c, query2.NewGetGroupBalance)
//	di.Provide(c, query4.NewGetGroupMonthlyIncome)
//	di.Provide(c, controller3.NewInviteUserToGroup)
//	di.Provide(c, controller3.NewAcceptGroupInvite)
//	di.Provide(c, controller3.NewGetGroupBalance)
//	di.Provide(c, controller5.NewGetGroupExpenses)
//	di.Provide(c, controller4.NewGetGroupMonthlyIncome)
//	di.Provide(c, controller3.NewCreateGroup)
//	di.Provide(c, controller3.NewGetGroup)
//	// expense
//	di.Provide(c, postgres5.NewExpenseRepository)
//	di.Provide(c, usecase6.NewCreateExpense)
//	di.Provide(c, usecase6.NewUpdateExpense)
//	di.Provide(c, usecase6.NewDeleteExpense)
//	di.Provide(c, usecase6.NewRecalculateExpensesSplitRatio)
//	di.Provide(c, query5.NewGetExpensesPerSearch)
//	di.Provide(c, query5.NewGetExpenseDetails)
//	di.Provide(c, controller5.NewGetExpensesPerSearch)
//	di.Provide(c, controller5.NewCreateExpense)
//	di.Provide(c, controller5.NewUpdateExpense)
//	di.Provide(c, controller5.NewDeleteExpense)
//	di.Provide(c, controller5.NewGetExpenseDetails)
//
//	// Insights
//	di.Provide(c, query5.NewGetExpensesPerPeriod)
//	di.Provide(c, query5.NewGetExpensesPerCategory)
//	di.Provide(c, query4.NewGetIncomesPerPeriod)
//	di.Provide(c, controller5.NewGetExpensesPerPeriod)
//	di.Provide(c, controller5.NewGetExpensesPerCategory)
//	di.Provide(c, controller4.NewGetIncomesPerPeriod)
//})
