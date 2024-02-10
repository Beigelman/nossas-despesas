package api

import (
	"errors"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	slog.ErrorContext(ctx.Context(), fmt.Sprintf("Error calling %s%s: erro = %s", ctx.BaseURL(), ctx.Path(), err.Error()), "error", err.Error())

	code := http.StatusInternalServerError
	message := http.StatusText(code)

	var e *except.HTTPError
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message.(string)
	}

	ctx.Set("Content-Type", "\"text/plain; charset=utf-8\"")

	return ctx.Status(code).JSON(ErrorResponse{
		StatusCode: code,
		Message:    message,
	})
}
