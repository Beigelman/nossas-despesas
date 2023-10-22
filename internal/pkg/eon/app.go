// Package eon is a micro framework that aims to make the processe of bootstrapping a new application quick and easy.The Tino APP will provide you a set of two important tools:
// - A dependency injection container
// - A life cycle manager for your application
package eon

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type applicationImpl struct {
	ctnr         *di.Container
	shutdownTime time.Duration
	lfcm         lifeCycleManager
	logger       Logger
	ctx          context.Context
	serviceName  string
}

func (app *applicationImpl) container() *di.Container {
	return app.ctnr
}

func (app *applicationImpl) BootStrap(modules ...Module) Application {
	var bootOrder []HookFn
	for i := range modules {
		module := modules[i]
		hookFn := func() error {
			app.logger.Info(fmt.Sprintf("Booting %s", module.name))
			module.bootFn(
				app.ctx,
				app.ctnr,
				app.lfcm,
				Info{
					ServiceName: app.serviceName,
				},
			)
			return nil
		}
		bootOrder = append(bootOrder, hookFn)
	}

	app.lfcm.OnBooting(HookOrders.APPEND, bootOrder...)

	return app
}

func (app *applicationImpl) Start() error {
	if err := app.lfcm.start(); err != nil {
		app.logger.Error("Error starting the application", "err", err)
		return fmt.Errorf("starting the application: %w", err)
	}

	signals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
		os.Kill,
	}

	ctx, stop := signal.NotifyContext(app.ctx, signals...)
	defer stop()
	// Waits for stop signal
	<-ctx.Done()
	app.lfcm.shutdown()

	return nil
}

func (app *applicationImpl) StartTest() error {
	if err := app.lfcm.start(); err != nil {
		app.logger.Error("Error starting the application", "err", err)
		return err
	}

	return nil
}

func (app *applicationImpl) Stop() error {
	if err := app.lfcm.stop(); err != nil {
		app.logger.Error("Failed to stop the application", "err", err)
		return fmt.Errorf("stopping the application: %w", err)
	}
	return nil
}
