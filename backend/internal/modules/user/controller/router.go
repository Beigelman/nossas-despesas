package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/shared/middleware"
)

func Router(
	server *fiber.App,
	getMyUserHandler GetMe,
	authMiddleware middleware.AuthMiddleware,
) {
	// Api group
	api := server.Group("api")
	// Api version V1
	v1 := api.Group("v1")
	user := v1.Group("user", authMiddleware)
	user.Get("/me", getMyUserHandler)
}
