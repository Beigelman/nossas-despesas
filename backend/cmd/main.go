package main

import (
	"log"
	"log/slog"
	"os"

	nossasdespesas "github.com/Beigelman/nossas-despesas"
	auth "github.com/Beigelman/nossas-despesas/internal/modules/auth/module"
	category "github.com/Beigelman/nossas-despesas/internal/modules/category/module"
	expense "github.com/Beigelman/nossas-despesas/internal/modules/expense/module"
	group "github.com/Beigelman/nossas-despesas/internal/modules/group/module"
	income "github.com/Beigelman/nossas-despesas/internal/modules/income/module"
	user "github.com/Beigelman/nossas-despesas/internal/modules/user/module"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/pkg/logger"
	"github.com/Beigelman/nossas-despesas/internal/shared"
)

func main() {
	environment := env.MustParse(os.Getenv("ENV"))

	cfg := nossasdespesas.MustNewConfig(environment)

	var lgr *slog.Logger
	if environment == env.Development {
		lgr = logger.NewDevelopment(cfg.LogLevel)
	} else {
		lgr = logger.NewProduction(cfg.LogLevel)
	}
	slog.SetDefault(lgr)

	ctnr := di.New()
	di.Concrete(ctnr, cfg)

	app := eon.New(cfg.ServiceName, eon.WithLogger(lgr), eon.WithIoC(ctnr))

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
