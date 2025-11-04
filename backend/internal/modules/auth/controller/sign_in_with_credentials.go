package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
)

type (
	SignInWithCredentialsRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	UserResponse struct {
		ID             int         `json:"id"`
		Name           string      `json:"name"`
		Email          string      `json:"email"`
		ProfilePicture *string     `json:"profile_picture,omitempty"`
		GroupID        *int        `json:"group_id,omitempty"`
		Flags          []user.Flag `json:"flags"`
		CreatedAt      time.Time   `json:"created_at"`
		UpdatedAt      time.Time   `json:"updated_at"`
	}

	UserLogIn struct {
		User         UserResponse `json:"user"`
		Token        string       `json:"token"`
		RefreshToken string       `json:"refresh_token"`
	}

	SignInWithCredentials func(ctx *fiber.Ctx) error
)

func NewSignInWithCredentials(signUpWithCredentials usecase.SignInWithCredentials) SignInWithCredentials {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req SignInWithCredentialsRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("ctx.BodyParser: %w", err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		result, err := signUpWithCredentials(ctx.Context(), usecase.SignInWithCredentialsParams{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			return fmt.Errorf("signUpWithCredentials: %w", err)
		}

		var groupID *int
		if result.User.GroupID != nil {
			groupID = &result.User.GroupID.Value
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, UserLogIn{
				User: UserResponse{
					ID:             result.User.ID.Value,
					Name:           result.User.Name,
					Email:          result.User.Email,
					ProfilePicture: result.User.ProfilePicture,
					GroupID:        groupID,
					Flags:          result.User.Flags,
					CreatedAt:      result.User.CreatedAt,
					UpdatedAt:      result.User.UpdatedAt,
				},
				Token:        result.Token,
				RefreshToken: result.RefreshToken,
			}),
		)
	}
}
