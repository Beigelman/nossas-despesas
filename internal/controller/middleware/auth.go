package middleware

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(auth *auth.Client) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		//tokenInfo := strings.Split(ctx.GetReqHeaders()["Authorization"], " ")
		//
		//if len(tokenInfo) != 2 {
		//	return except.UnauthorizedError("invalid token format")
		//}
		//
		//prefix := tokenInfo[0]
		//token := tokenInfo[1]
		//
		//if token == "" || prefix != "Bearer" {
		//	return except.UnauthorizedError("invalid token format")
		//}
		//
		//idToken, err := auth.VerifyIDToken(ctx.Context(), token)
		//if err != nil {
		//	return except.UnauthorizedError("invalid token").SetInternal(err)
		//}
		//
		//ctx.Set("userID", idToken.Claims["userID"].(string))
		//ctx.Set("email", idToken.Claims["email"].(string))
		//ctx.Set("groupID", idToken.Claims["groupID"].(string))

		return ctx.Next()
	}
}
