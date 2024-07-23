package boot

import (
	"context"
	"errors"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
	"github.com/Beigelman/nossas-despesas/internal/shared/middleware"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
	"log/slog"
	"net/http"
	"time"
)

var ServerModule = eon.NewModule("Server", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	var server *fiber.App

	di.Provide(c, middleware.NewAuthMiddleware)

	di.Provide(c, func(cfg *config.Config) *fiber.App {
		server = fiber.New(fiber.Config{
			AppName:      info.ServiceName,
			ReadTimeout:  5 * time.Second,
			JSONEncoder:  sonic.Marshal,
			JSONDecoder:  sonic.Unmarshal,
			ErrorHandler: api.ErrorHandler,
		})

		server.Use(cors.New())
		server.Use(recover.New())
		server.Use(requestid.New())
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
