package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/Beigelman/nossas-despesas/internal/pkg/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/shared/middleware"
)

var Module = eon.NewModule("Server", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	var server *fiber.App

	di.Provide(c, middleware.NewAuthMiddleware)

	di.Provide(c, func(cfg *config.Config) *fiber.App {
		server = fiber.New(fiber.Config{
			AppName:      info.ServiceName,
			ReadTimeout:  5 * time.Second,
			ErrorHandler: ErrorHandler,
		})

		server.Use(cors.New())
		server.Use(recover.New())
		server.Use(requestid.New())
		server.Use(sentryfiber.New(sentryfiber.Options{
			Repanic:         true,
			WaitForDelivery: true,
		}))
		server.Use(middleware.LogRequest(cfg.Env, info.ServiceName))

		server.Get("health-check", func(c *fiber.Ctx) error { return c.SendString("OK") })

		return server
	})

	lc.OnReady(eon.HookOrders.APPEND, func() error {
		go func() {
			cfg := di.Resolve[*config.Config](c)
			if err := server.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(fmt.Errorf("server.Listen: %w", err))
			}
		}()

		return nil
	})

	lc.OnDisposing(eon.HookOrders.APPEND, func() error {
		if server != nil {
			slog.Info("Shutting down server")
			if err := server.Shutdown(); err != nil {
				return fmt.Errorf("server.Shutdown: %w", err)
			}
		}
		return nil
	})
})
