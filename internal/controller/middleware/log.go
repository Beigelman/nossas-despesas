package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func LogRequest(serviceName string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("x-service-name", serviceName)
		slog.Info(
			fmt.Sprintf("Calling %s%s", ctx.BaseURL(), ctx.Path()),
			"method", ctx.Method(),
			"ip", ctx.IP(),
			"service", serviceName,
		)
		return ctx.Next()
	}
}
