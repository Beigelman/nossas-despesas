package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/Beigelman/nossas-despesas/internal/config"
	auth "github.com/Beigelman/nossas-despesas/internal/modules/auth/module"
	category "github.com/Beigelman/nossas-despesas/internal/modules/category/module"
	expense "github.com/Beigelman/nossas-despesas/internal/modules/expense/module"
	group "github.com/Beigelman/nossas-despesas/internal/modules/group/module"
	income "github.com/Beigelman/nossas-despesas/internal/modules/income/module"
	user "github.com/Beigelman/nossas-despesas/internal/modules/user/module"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/pkg/logger"
	"github.com/Beigelman/nossas-despesas/internal/shared"
)

func main() {
	environment := env.MustParse(os.Getenv("ENV"))

	var lgr *slog.Logger
	if environment == env.Development {
		lgr = logger.NewDevelopment()
	} else {
		lgr = logger.NewProduction()
	}

	app := eon.New("Nossas Despesas", eon.WithLogger(lgr))

	if err := app.BootStrap(
		// Common Modules
		config.Module,
		db.Module,
		api.Module,
		shared.Module,
		// Domain Modules
		auth.Module,
		category.Module,
		expense.Module,
		group.Module,
		income.Module,
		user.Module,
	).Start(); err != nil {
		log.Fatal("failed to start application: ", err)
	}
}
