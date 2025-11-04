package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	requestId, ok := ctx.Locals("requestid").(string)
	if !ok {
		requestId = "unknown"
	}

	code := http.StatusInternalServerError
	message := http.StatusText(code)
	errMsg := err.Error()
	var e *except.HTTPError
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message.(string)
	}

	hub := sentryfiber.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.WithScope(func(scope *sentry.Scope) {
		scope.AddEventProcessor(func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			hint.Context = ctx.Context()
			for _, e := range event.Exception {
				if e.Stacktrace != nil {
					e.Stacktrace.Frames = e.Stacktrace.Frames[:len(e.Stacktrace.Frames)-2]
				}
			}
			return event
		})

		scope.SetTag("request_id", requestId)
		scope.SetContext("error", sentry.Context{
			"code":    code,
			"message": message,
			"error":   errMsg,
			"method":  ctx.Method(),
			"path":    ctx.Path(),
		})

		hub.CaptureException(err)
	})

	slog.Error(
		fmt.Sprintf("Error calling %s %s", ctx.Method(), ctx.Path()),
		slog.String("request_id", requestId),
		slog.Int("status_code", code),
		slog.String("error", errMsg),
	)

	ctx.Set("Content-Type", "\"text/plain; charset=utf-8\"")

	return ctx.Status(code).JSON(ErrorResponse{
		StatusCode: code,
		Message:    message,
		Error:      errMsg,
	})
}
