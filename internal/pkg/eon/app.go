// Eon APP is a micro framework that aims to make the processe of bootstrapping a new application quick and easy.The Tino APP will provide you a set of two important tools:
// - A dependency injection container
// - A life cycle manager for your application
package eon

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
	"github.com/truepay/go-commons/logger"
	"github.com/truepay/go-commons/logger/field"
)

type Application interface {
	// The BootStrap method is used to register all the modules that your application will need in the booting process.
	BootStrap(modules ...Module) Application
	// The Start method is used to start the application and get through the firsts lifecycle hooks.
	// This method will block the current goroutine execution.
	Start() error
	// The StartTest method is used to start the application without blocking the execution of the current goroutine.
	StartTest() error
	// The Stop method is used to stop the application and get through the last lifecycle hooks.
	Stop() error
	// The Container method is used to get the dependency injection container.
	container() *di.Container
}

type applicationImpl struct {
	ctnr        *di.Container
	lfcm        lifeCycleManager
	logger      *logger.CtxLogger
	env         env.Environment
	ctx         context.Context
	serviceName string
}

type Config struct {
	ShutdownTime time.Duration
	ServiceName  string
	Env          env.Environment
	IoC          *di.Container
}

// New creates a new concrete of the eon Application.
// Eon APP is a micro framework that aims to make the processe of bootstrapping a new application quick and easy.
// The Tino APP will provide you a set of two important tools: a dependency injection container and a life cycle manager for your application
func NewApp(cfg Config) Application {
	ctx := logger.ContextWithArgs(context.Background(), field.String("service", cfg.ServiceName), field.String("env", string(cfg.Env)))

	var (
		tlog *logger.Logger
		err  error
	)
	if cfg.Env == env.Production {
		tlog, err = logger.NewProductionLogger(cfg.ServiceName, logger.InfoLevel)
	} else {
		tlog, err = logger.NewDevelopmentLogger(cfg.ServiceName, logger.InfoLevel)
	}
	if err != nil {
		log.Fatal(err)
	}

	var ctnr *di.Container
	if cfg.IoC != nil {
		ctnr = cfg.IoC
	} else {
		ctnr = di.New()
	}

	return &applicationImpl{
		serviceName: cfg.ServiceName,
		env:         cfg.Env,
		ctnr:        ctnr,
		ctx:         ctx,
		logger:      tlog.Ctx(ctx),
		lfcm:        newLifeCycleManager(cfg.ShutdownTime, tlog, ctx),
	}
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
					Env:         app.env,
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
	if err := app.lfcm.start(app.ctx); err != nil {
		app.logger.Error("Error starting the application", field.Error(err))
		return err
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
	app.lfcm.shutdown(app.ctx)

	return nil
}

func (app *applicationImpl) StartTest() error {
	if err := app.lfcm.start(app.ctx); err != nil {
		app.logger.Error("Error starting the application", field.Error(err))
		return err
	}
	return nil
}

func (app *applicationImpl) Stop() error {
	if err := app.lfcm.stop(app.ctx); err != nil {
		app.logger.Error("Failed to stop the application", field.Error(err))
		return err
	}
	return nil
}
