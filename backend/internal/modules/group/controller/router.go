package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/shared/middleware"
)

func Router(
	server *fiber.App,
	createGroupHandler CreateGroup,
	getGroupHandler GetGroup,
	authMiddleware middleware.AuthMiddleware,
	inviteUserToGroupHandler InviteUserToGroup,
	acceptGroupInviteHandler AcceptGroupInvite,
	getGroupBalanceHandler GetGroupBalance,
) {
	// Api group
	api := server.Group("api")
	// Api version V1
	v1 := api.Group("v1")
	// Group routes
	group := v1.Group("group", authMiddleware)
	group.Get("/", getGroupHandler)
	group.Post("/", createGroupHandler)
	group.Get("/balance", getGroupBalanceHandler)
	// Invite Router
	invite := group.Group("invite", authMiddleware)
	invite.Post("/", inviteUserToGroupHandler)
	invite.Post("/:token/accept", acceptGroupInviteHandler)
}
