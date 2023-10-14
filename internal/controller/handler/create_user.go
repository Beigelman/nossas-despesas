package handler

import (
	"fmt"
	"net/http"

	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type (
	CreateUserRequest struct {
		Name           string  `json:"name"`
		Email          string  `json:"email"`
		ProfilePicture *string `json:"profile_picture"`
	}
	CreateUserResponse struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	CreateUser func(ctx *fiber.Ctx) error
)

func NewCreateUserHandler(createUser usecase.CreateUser) CreateUser {
	return func(ctx *fiber.Ctx) error {
		var req CreateUserRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("ctx.BodyParser: %w", err)
		}

		user, err := createUser(ctx.Context(), usecase.CreateUserParams{
			Name:           req.Name,
			Email:          req.Email,
			ProfilePicture: req.ProfilePicture,
		})
		if err != nil {
			return fmt.Errorf("createUser: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, CreateUserResponse{
				ID:    user.ID.Value,
				Name:  user.Name,
				Email: user.Email,
			}),
		)
	}
}
