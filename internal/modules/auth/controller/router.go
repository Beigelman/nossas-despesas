package controller

import (
	"github.com/gofiber/fiber/v2"
)

func Router(
	server *fiber.App,
	signInWithCredentialsHandler SignInWithCredentials,
	signUpWithCredentialsHandler SignUpWithCredentials,
	signInWithGoogleHandler SignInWithGoogle,
	refreshAuthTokenHandler RefreshAuthToken,

) {

	// Api group
	api := server.Group("api")
	// Api version V1

	v1 := api.Group("v1")
	// Auth routes
	auth := v1.Group("auth")
	auth.Post("/sign-in/credentials", signInWithCredentialsHandler)
	auth.Post("/sign-in/google", signInWithGoogleHandler)
	auth.Post("/sign-up/credentials", signUpWithCredentialsHandler)
	auth.Post("refresh-token", refreshAuthTokenHandler)
}
