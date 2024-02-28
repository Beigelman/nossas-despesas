package main

import (
	"github.com/Beigelman/nossas-despesas/internal/boot"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/pkg/logger"
	"log"
	"log/slog"
	"os"
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
		boot.ConfigModule,
		boot.ClientsModule,
		boot.ServerModule,
		boot.ApplicationModule,
	).Start(); err != nil {
		log.Fatal("failed to start application: ", err)
	}
}
