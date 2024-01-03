package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/pkg/validator"
	"net/http"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type (
	CreateUserRequest struct {
		Name             string  `json:"name" validate:"required"`
		Email            string  `json:"email" validate:"required"`
		ProfilePicture   *string `json:"profile_picture"`
		AuthenticationID *string `json:"authentication_id"`
		GroupID          *int    `json:"group_id"`
	}

	CreateUserResponse struct {
		ID               int       `json:"id"`
		Name             string    `json:"name"`
		Email            string    `json:"email"`
		GroupID          int       `json:"group_id,omitempty"`
		AuthenticationID string    `json:"authentication_id,omitempty"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
	}

	CreateUser func(ctx *fiber.Ctx) error
)

func NewCreateUser(createUser usecase.CreateUser) CreateUser {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req CreateUserRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("ctx.BodyParser: %w", err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		user, err := createUser(ctx.Context(), usecase.CreateUserParams{
			Name:             req.Name,
			Email:            req.Email,
			ProfilePicture:   req.ProfilePicture,
			AuthenticationID: req.AuthenticationID,
			GroupID: func() *entity.GroupID {
				if req.GroupID == nil {
					return nil
				}
				return &entity.GroupID{Value: *req.GroupID}
			}(),
		})
		if err != nil {
			return fmt.Errorf("createUser: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, CreateUserResponse{
				ID:    user.ID.Value,
				Name:  user.Name,
				Email: user.Email,
				GroupID: func() int {
					if user.GroupID == nil {
						return 0
					}
					return user.GroupID.Value
				}(),
				AuthenticationID: func() string {
					if user.AuthenticationID == nil {
						return ""
					}
					return *user.AuthenticationID
				}(),
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}),
		)
	}
}
