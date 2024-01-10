package middleware

import (
	"github.com/Beigelman/ludaapi/internal/infra/jwt"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

type AuthMiddleware func(ctx *fiber.Ctx) error

func NewAuthMiddleware(tokenProvider *jwt.Provider) AuthMiddleware {
	return func(ctx *fiber.Ctx) error {
		tokenInfo := strings.Split(ctx.GetReqHeaders()["Authorization"], " ")

		if len(tokenInfo) != 2 {
			return except.UnauthorizedError("invalid jwt format")
		}

		prefix := tokenInfo[0]
		token := tokenInfo[1]

		if token == "" || prefix != "Bearer" {
			return except.UnauthorizedError("invalid jwt format")
		}

		idToken, err := tokenProvider.ParseToken(token)
		if err != nil {
			return except.UnauthorizedError("invalid jwt").SetInternal(err)
		}

		claims, ok := idToken.Claims.(jwt.MapClaims)
		if !ok {
			return except.UnauthorizedError("invalid jwt")
		}

		ctx.Set("userID", claims["userID"].(string))
		ctx.Set("email", claims["email"].(string))
		ctx.Set("groupID", claims["groupID"].(string))

		return ctx.Next()
	}
}
