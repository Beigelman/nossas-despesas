package main

import (
	"log"
	"os"
	"time"

	"github.com/Beigelman/ludaapi/internal/boot"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
	"github.com/Beigelman/ludaapi/internal/pkg/eon"
)

func main() {
	environment, err := env.Parse(os.Getenv("ENV"))
	if err != nil {
		log.Fatal("failed to parse environment: ", err)
	}

	app := eon.NewApp(eon.Config{
		ServiceName:  "Luda API",
		ShutdownTime: 10 * time.Second,
		Env:          environment,
	})

	if err := app.BootStrap(boot.ConfigModule, boot.DatabaseModule, boot.ServerModule, boot.ApplicationModule).Start(); err != nil {
		log.Fatal("failed to start application: ", err)
	}
}
