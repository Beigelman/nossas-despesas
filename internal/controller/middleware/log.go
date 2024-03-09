package middleware

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/gofiber/fiber/v2"
)

func LogRequest(environment env.Environment, serviceName string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		if environment == "test" {
			return ctx.Next()
		}

		requestId, ok := ctx.Locals("requestid").(string)
		if !ok {
			panic("request_id not found in context")
		}

		if environment == "development" {
			slog.Info(fmt.Sprintf("Calling %s %s", ctx.Method(), ctx.Path()), slog.String("request_id", requestId))
			return ctx.Next()
		}

		params := extractRequestParams(ctx)
		slog.Info(fmt.Sprintf("Calling %s %s", ctx.Method(), ctx.Path()),
			slog.String("request_id", requestId),
			slog.Group("http_request",
				slog.String("ip", ctx.IP()),
				slog.String("service", serviceName),
				slog.String("method", ctx.Method()),
				slog.Group("headers", params.headerParams...),
				slog.Group("query", params.queryParams...),
				slog.Group("body", params.bodyParams...),
			))
		return ctx.Next()
	}
}

var secretRegex = regexp.MustCompile(`(?i)password|token|secret|key|authorization|session|jwt|auth`)

type requestParams struct {
	bodyParams   []any
	queryParams  []any
	headerParams []any
}

func extractRequestParams(ctx *fiber.Ctx) requestParams {
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
			value := fmt.Sprintf("%s", v[0])
			if secretRegex.MatchString(k) {
				value = "***"
			}
			headers = append(headers, slog.String(k, value))
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
			value := fmt.Sprintf("%s", v)
			if secretRegex.MatchString(k) {
				value = "***"
			}
			bodyParams = append(bodyParams, slog.String(k, value))
		}

		return bodyParams
	}()

	return requestParams{
		bodyParams:   bodyParams,
		queryParams:  queryParams,
		headerParams: headerParams,
	}
}
