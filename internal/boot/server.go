package boot

import (
	"context"
	"errors"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/controller"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/Beigelman/ludaapi/internal/pkg/eon"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

var ServerModule = eon.NewModule("Server", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	var server *fiber.App

	di.Provide(c, func() *fiber.App {
		server = fiber.New(fiber.Config{
			AppName:      info.ServiceName,
			ReadTimeout:  5 * time.Second,
			JSONEncoder:  sonic.Marshal,
			JSONDecoder:  sonic.Unmarshal,
			ErrorHandler: errorHandler,
		})

		server.Use(func(ctx *fiber.Ctx) error {
			ctx.Set("x-service-name", info.ServiceName)
			slog.Info(fmt.Sprintf("Calling %s%s", ctx.BaseURL(), ctx.Path()), "method", ctx.Method(), "ip", ctx.IP())
			return ctx.Next()
		})

		return server
	})

	lc.OnBooted(eon.HookOrders.PREPEND, func() error {
		return di.Call(c, controller.Router)
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
		slog.Info("Shutting down server")
		if err := server.Shutdown(); err != nil {
			return fmt.Errorf("server.Shutdown: %w", err)
		}
		return nil
	})
})

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := http.StatusInternalServerError
	message := http.StatusText(code)
	var e *except.HTTPError
	if errors.As(err, &e) {
		code = e.Code
		message = e.Error()
	}
	ctx.Set("Content-Type", "\"text/plain; charset=utf-8\"")
	return ctx.Status(code).JSON(map[string]any{
		"code":    code,
		"message": message,
	})
}
