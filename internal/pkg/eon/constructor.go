package eon

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
)

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Debug(msg string, fields ...any)
}

type Options func(app *Application)

// New creates a new concrete of the eon Application.
// Eon APP is a micro framework that aims to make the processe of bootstrapping a new application quick and easy.
// The Tino APP will provide you a set of two important tools: a dependency injection container and a life cycle manager for your application
func New(serviceName string, opts ...Options) *Application {
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	app := &Application{
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
	return func(app *Application) {
		app.logger = logger
		app.lfcm = newLifeCycleManager(app.shutdownTime, logger)
	}
}

func WithShutdownTime(shutdownTime time.Duration) Options {
	return func(app *Application) {
		app.shutdownTime = shutdownTime
		app.lfcm = newLifeCycleManager(shutdownTime, app.logger)
	}
}

func WithIoC(ioc *di.Container) Options {
	return func(app *Application) {
		app.ctnr = ioc
	}
}
