package middleware

import (
	"github.com/Beigelman/ludaapi/internal/domain/service"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type AuthMiddleware func(ctx *fiber.Ctx) error

func NewAuthMiddleware(tokenProvider service.TokenProvider) AuthMiddleware {
	return func(ctx *fiber.Ctx) error {
		bearerToken := strings.Split(ctx.GetReqHeaders()["Authorization"], " ")

		if len(bearerToken) != 2 {
			return except.UnauthorizedError("invalid jwt format")
		}

		prefix := bearerToken[0]
		token := bearerToken[1]

		if token == "" || prefix != "Bearer" {
			return except.UnauthorizedError("invalid jwt format")
		}

		tokenInfo, err := tokenProvider.ParseToken(token)
		if err != nil {
			return except.UnauthorizedError("invalid jwt").SetInternal(err)
		}

		ctx.Locals("user_id", tokenInfo.Claims.UserID)
		ctx.Locals("email", tokenInfo.Claims.Email)
		if tokenInfo.Claims.GroupID != nil {
			ctx.Locals("group_id", *tokenInfo.Claims.GroupID)
		}

		return ctx.Next()
	}
}
