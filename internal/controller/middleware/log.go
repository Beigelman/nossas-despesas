package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"regexp"
)

var secretRegex = regexp.MustCompile(`password|token|secret|key`)

func LogRequest(serviceName string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		queryParams := func() []any {
			var query []any
			for k, v := range ctx.Queries() {
				query = append(query, slog.String(k, v))
			}
			return query
		}()

		headerParams := func() []any {
			var headers []any
			for k, v := range ctx.GetReqHeaders() {
				headers = append(headers, slog.String(k, v[0]))
			}
			return headers
		}()

		bodyParams := func() []any {
			var body map[string]any
			if err := ctx.BodyParser(&body); err != nil {
				return nil
			}

			var bodyParams []any
			for k, v := range body {
				value := fmt.Sprintf("%v", v)
				if secretRegex.MatchString(k) {
					value = "***"
				}
				bodyParams = append(bodyParams, slog.String(k, value))
			}

			return bodyParams
		}()

		slog.Info(fmt.Sprintf("Calling %s", ctx.Path()),
			slog.String("request_id", ctx.Locals("requestid").(string)),
			slog.Group("http_request",
				slog.String("ip", ctx.IP()),
				slog.String("service", serviceName),
				slog.String("method", ctx.Method()),
				slog.Group("headers", headerParams...),
				slog.Group("query", queryParams...),
				slog.Group("body", bodyParams...),
			))
		return ctx.Next()
	}
}
