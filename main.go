package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/Beigelman/nossas-despesas/internal/boot"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/pkg/logger"
)

func main() {
	environment, err := env.Parse(os.Getenv("ENV"))
	if err != nil {
		log.Fatal("failed to parse environment: ", err)
	}

	var lgr *slog.Logger
	if environment == env.Development {
		lgr = logger.NewDevelopment()
	} else {
		lgr = logger.NewProduction()
	}

	app := eon.New("Nossas Despesas", eon.WithLogger(lgr))

	if err := app.BootStrap(
		// Common Modules
		boot.ConfigModule,
		boot.DatabaseModule,
		boot.ServerModule,
		boot.ClientsModule,
		// Domain Modules
		boot.AuthModule,
		boot.CategoryModule,
		boot.ExpenseModule,
		boot.GroupModule,
		boot.IncomeModule,
		boot.UserModule,
	).Start(); err != nil {
		log.Fatal("failed to start application: ", err)
	}
}
