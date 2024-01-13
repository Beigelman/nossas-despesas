package middleware

import (
	"fmt"
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

		ctx.Set("userID", fmt.Sprintf("%d", tokenInfo.Claims.UserID))
		ctx.Set("groupID", fmt.Sprintf("%d", tokenInfo.Claims.UserID))
		ctx.Set("email", tokenInfo.Claims.Email)

		return ctx.Next()
	}
}
