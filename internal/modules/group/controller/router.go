package controller

import (
	"github.com/Beigelman/nossas-despesas/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(
	server *fiber.App,
	createGroupHandler CreateGroup,

	getGroupHandler GetGroup,

	authMiddleware middleware.AuthMiddleware,
	inviteUserToGroupHandler InviteUserToGroup,
	acceptGroupInviteHandler AcceptGroupInvite,
) {
	// Api group
	api := server.Group("api")
	// Api version V1
	v1 := api.Group("v1")
	// Group routes
	group := v1.Group("group", authMiddleware)
	group.Get("/", getGroupHandler)
	group.Post("/", createGroupHandler)
	// Invite Router
	group.Post("/invite", inviteUserToGroupHandler)
	group.Post("/invite/:token/accept", acceptGroupInviteHandler)
}
