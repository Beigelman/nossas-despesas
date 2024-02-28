package main

import (
	"github.com/Beigelman/nossas-despesas/internal/boot"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"log"
)

func main() {
	app := eon.New("Luda API")

	if err := app.BootStrap(
		boot.ConfigModule,
		boot.ClientsModule,
		boot.ServerModule,
		boot.ApplicationModule,
	).Start(); err != nil {
		log.Fatal("failed to start application: ", err)
	}
}
