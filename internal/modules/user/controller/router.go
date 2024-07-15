package controller

import (
	"github.com/Beigelman/nossas-despesas/internal/modules/user/controller/handler"
	"github.com/Beigelman/nossas-despesas/internal/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(
	server *fiber.App,
	getMyUserHandler handler.GetMe,
	authMiddleware middleware.AuthMiddleware,
) {
	// Api group
	api := server.Group("api")
	// Api version V1
	v1 := api.Group("v1")
	user := v1.Group("user", authMiddleware)
	user.Get("/me", getMyUserHandler)

}
