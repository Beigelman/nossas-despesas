package api

import (
	"errors"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	slog.Error(
		fmt.Sprintf("Error calling %s %s", ctx.Method(), ctx.Path()),
		slog.String("request_id", ctx.Locals("requestid").(string)),
		slog.String("error", err.Error()),
	)

	code := http.StatusInternalServerError
	message := http.StatusText(code)
	errMsg := ""
	var e *except.HTTPError
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message.(string)
		errMsg = e.Error()
	}

	ctx.Set("Content-Type", "\"text/plain; charset=utf-8\"")

	return ctx.Status(code).JSON(ErrorResponse{
		StatusCode: code,
		Message:    message,
		Error:      errMsg,
	})
}
