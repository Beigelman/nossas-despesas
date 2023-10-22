package main

import (
	"github.com/Beigelman/ludaapi/internal/boot"
	"github.com/Beigelman/ludaapi/internal/pkg/eon"
	"log"
)

func main() {
	app := eon.New("Luda API")

	if err := app.BootStrap(
		boot.ConfigModule,
		boot.DatabaseModule,
		boot.ServerModule,
		boot.ApplicationModule,
	).Start(); err != nil {
		log.Fatal("failed to start application: ", err)
	}
}
