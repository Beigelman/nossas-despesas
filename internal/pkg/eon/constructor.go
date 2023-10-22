package eon

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"log/slog"
	"os"
	"time"
)

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Debug(msg string, fields ...any)
}

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

type Options func(app *applicationImpl)

// New creates a new concrete of the eon Application.
// Eon APP is a micro framework that aims to make the processe of bootstrapping a new application quick and easy.
// The Tino APP will provide you a set of two important tools: a dependency injection container and a life cycle manager for your application
func New(serviceName string, opts ...Options) Application {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).WithGroup("eon")
	app := &applicationImpl{
		serviceName:  serviceName,
		ctnr:         di.New(),
		ctx:          context.Background(),
		logger:       l,
		shutdownTime: 10 * time.Second,
		lfcm:         newLifeCycleManager(10*time.Second, l),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func WithLogger(logger Logger) Options {
	return func(app *applicationImpl) {
		app.logger = logger
		app.lfcm = newLifeCycleManager(app.shutdownTime, logger)
	}
}

func WithShutdownTime(shutdownTime time.Duration) Options {
	return func(app *applicationImpl) {
		app.shutdownTime = shutdownTime
		app.lfcm = newLifeCycleManager(shutdownTime, app.logger)
	}
}

func WithIoC(ioc *di.Container) Options {
	return func(app *applicationImpl) {
		app.ctnr = ioc
	}
}
